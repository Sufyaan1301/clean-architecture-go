package periksa

import (
	"be-golang-poliklinik/internal/domain"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DetailObat struct {
	IDObat int
	Jumlah int
}

type Repository interface {
	FindLastAntrian(idJadwal int) int
	SavePeriksaByPasien(daftar *domain.DaftarPoli) error
	FindPeriksaByDokter(idDokter int, hariIni string) ([]domain.DaftarPoli, error)
	SavePeriksaByDokterTx(periksa *domain.Periksa, obatInputs []DetailObat) error
	FindLastPeriksaByPasien(idPasien int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindLastAntrian(idJadwal int) int {
	var daftar domain.DaftarPoli
	err := r.db.Where("id_jadwal = ?", idJadwal).Order("no_antrian desc").First(&daftar).Error
	if err != nil {
		return 0
	}
	return daftar.NoAntrian
}

func (r *repository) SavePeriksaByPasien(daftar *domain.DaftarPoli) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var lastDaftar domain.DaftarPoli
		
		// PESSIMISTIC LOCKING: Kunci data antrian jadwal ini selama transaksi berlangsung
		// Ini mencegah goroutine/request lain membaca 'no_antrian' yang sama.
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id_jadwal = ?", daftar.IDJadwal).
			Order("no_antrian desc").
			First(&lastDaftar).Error
		
		noAntrian := 1
		if err == nil {
			noAntrian = lastDaftar.NoAntrian + 1
		}
		
		daftar.NoAntrian = noAntrian
		return tx.Create(daftar).Error
	})

	if err == nil {
		r.db.Preload("Pasien").Preload("Jadwal").Preload("Jadwal.Dokter").First(daftar, daftar.ID)
	}
	return err
}

func (r *repository) FindPeriksaByDokter(idDokter int, hariIni string) ([]domain.DaftarPoli, error) {
	var antrian []domain.DaftarPoli
	err := r.db.
		Preload("Pasien").
		Preload("Jadwal").
		Joins("JOIN jadwal_periksas ON jadwal_periksas.id = daftar_polis.id_jadwal").
		Joins("LEFT JOIN periksas ON periksas.id_daftar_poli = daftar_polis.id").
		Where("jadwal_periksas.id_dokter = ?", idDokter).
		Where("jadwal_periksas.hari = ?", hariIni).
		Where("periksas.id IS NULL").
		Find(&antrian).Error
	return antrian, err
}

func (r *repository) SavePeriksaByDokterTx(periksa *domain.Periksa, obatInputs []DetailObat) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		totalBiaya := 150000 // Biaya dasar jasa dokter
		if err := tx.Create(periksa).Error; err != nil {
			return err
		}
		for _, input := range obatInputs {
			var obat domain.Obat
			if err := tx.First(&obat, input.IDObat).Error; err != nil {
				return errors.New("terdapat obat yang tidak ditemukan")
			}

			if obat.Stok < input.Jumlah {
				return errors.New("stok obat " + obat.NamaObat + " tidak mencukupi")
			}
			
			// Kurangi stok sesuai jumlah yang diminta
			obat.Stok -= input.Jumlah
			if err := tx.Save(&obat).Error; err != nil {
				return err
			}
			
			// Tambahkan harga obat dikali jumlahnya
			totalBiaya += (obat.Harga * input.Jumlah)
			
			detail := domain.DetailPeriksa{
				IDPeriksa: periksa.ID,
				IDObat:    obat.ID,
				Jumlah:    input.Jumlah,
			}
			if err := tx.Create(&detail).Error; err != nil {
				return err
			}
		}
		periksa.BiayaPeriksa = totalBiaya
		return tx.Save(periksa).Error
	})
}

func (r *repository) FindLastPeriksaByPasien(idPasien int) error {
	// CEK 1: Apakah pasien masih antri (belum diperiksa)?
	var countAntrian int64
	r.db.Table("daftar_polis").
		Joins("LEFT JOIN periksas ON periksas.id_daftar_poli = daftar_polis.id").
		Where("daftar_polis.id_pasien = ?", idPasien).
		Where("periksas.id IS NULL").
		Count(&countAntrian)

	if countAntrian > 0 {
		return errors.New("Anda masih ada Booking pemeriksaan yang aktif")
	}

	// CEK 2: Apakah pasien sudah diperiksa tapi BELUM BAYAR?
	var countTagihan int64
	r.db.Table("periksas").
		Joins("JOIN daftar_polis ON daftar_polis.id = periksas.id_daftar_poli").
		Joins("LEFT JOIN pembayarans ON pembayarans.id_periksa = periksas.id").
		Where("daftar_polis.id_pasien = ?", idPasien).
		Where("pembayarans.id IS NULL OR pembayarans.status != 'Lunas'"). // Harus benar-benar Lunas
		Count(&countTagihan)

	if countTagihan > 0 {
		return errors.New("Pendaftaran ditolak: Anda memiliki tagihan pemeriksaan sebelumnya yang belum berstatus Lunas.")
	}

	return nil
}
