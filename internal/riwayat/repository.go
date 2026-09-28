package riwayat

import (
	"be-golang-poliklinik/internal/domain"

	"gorm.io/gorm"
)

type Repository interface {
	GetRiwayatDokter(idDokter int) ([]domain.Periksa, error)
	GetRiwayatDokterByID(idPeriksa int, idDokter int) (*domain.Periksa, error)
	GetRiwayatPasien(idPasien int) ([]domain.Periksa, error)
	GetRiwayatPasienByID(idPeriksa int, idPasien int) (*domain.Periksa, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetRiwayatDokter(idDokter int) ([]domain.Periksa, error) {
	var riwayats []domain.Periksa
	err := r.db.
		Preload("DaftarPoli").
		Preload("DaftarPoli.Pasien").
		Preload("DetailPeriksa").
		Preload("DetailPeriksa.Obat").
		Joins("JOIN daftar_polis ON daftar_polis.id = periksas.id_daftar_poli").
		Joins("JOIN jadwal_periksas ON jadwal_periksas.id = daftar_polis.id_jadwal").
		Where("jadwal_periksas.id_dokter = ?", idDokter).
		Order("periksas.tgl_periksa DESC").
		Find(&riwayats).Error
	return riwayats, err
}

func (r *repository) GetRiwayatPasien(idPasien int) ([]domain.Periksa, error) {
	var riwayats []domain.Periksa
	err := r.db.
		Preload("DaftarPoli").
		Preload("DaftarPoli.Jadwal.Dokter.Poli").
		Preload("DetailPeriksa").
		Preload("DetailPeriksa.Obat").
		Preload("Pembayaran"). // Tarik status bayarnya!
		Joins("JOIN daftar_polis ON daftar_polis.id = periksas.id_daftar_poli").
		Where("daftar_polis.id_pasien = ?", idPasien).
		Order("periksas.tgl_periksa DESC").
		Find(&riwayats).Error
	return riwayats, err
}

func (r *repository) GetRiwayatDokterByID(idPeriksa int, idDokter int) (*domain.Periksa, error) {
	var riwayat domain.Periksa
	err := r.db.
		Preload("DaftarPoli").
		Preload("DaftarPoli.Pasien").
		Preload("DetailPeriksa").
		Preload("DetailPeriksa.Obat").
		Joins("JOIN daftar_polis ON daftar_polis.id = periksas.id_daftar_poli").
		Joins("JOIN jadwal_periksas ON jadwal_periksas.id = daftar_polis.id_jadwal").
		Where("jadwal_periksas.id_dokter = ?", idDokter).
		Where("periksas.id = ?", idPeriksa).
		First(&riwayat).Error
	return &riwayat, err
}

func (r *repository) GetRiwayatPasienByID(idPeriksa int, idPasien int) (*domain.Periksa, error) {
	var riwayat domain.Periksa
	err := r.db.
		Preload("DaftarPoli").
		Preload("DaftarPoli.Jadwal.Dokter.Poli").
		Preload("DetailPeriksa").
		Preload("DetailPeriksa.Obat").
		Preload("Pembayaran").
		Joins("JOIN daftar_polis ON daftar_polis.id = periksas.id_daftar_poli").
		Where("daftar_polis.id_pasien = ?", idPasien).
		Where("periksas.id = ?", idPeriksa).
		First(&riwayat).Error
	return &riwayat, err
}