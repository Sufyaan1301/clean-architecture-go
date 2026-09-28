package periksa

import (
	"be-golang-poliklinik/internal/domain"
	"errors"
	"time"
)

type Service interface {
	CreateDaftarPoli(idPasien int, input DaftarPoliInput) (*domain.DaftarPoli, error)
	GetAntrianDokter(idDokter int, hariIni string) ([]domain.DaftarPoli, error)
	CreatePeriksa(input PeriksaInput) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

type DaftarPoliInput struct {
	IDJadwal int    `json:"id_jadwal" binding:"required"`
	Keluhan  string `json:"keluhan" binding:"required"`
}

type ObatInput struct {
	IDObat int `json:"id_obat" binding:"required"`
	Jumlah int `json:"jumlah" binding:"required,min=1"`
}

type PeriksaInput struct {
	IDDaftarPoli int         `json:"id_daftar_poli" binding:"required"`
	Catatan      string      `json:"catatan" binding:"required"`
	Obat         []ObatInput `json:"obat" binding:"required,dive"`
}

func (s *service) CreateDaftarPoli(idPasien int, input DaftarPoliInput) (*domain.DaftarPoli, error) {

	hasTanggungan := s.repo.FindLastPeriksaByPasien(idPasien)
	if hasTanggungan != nil {
		return nil, hasTanggungan
	}

	lastAntrian := s.repo.FindLastAntrian(input.IDJadwal)

	daftar := domain.DaftarPoli{
		IDPasien:  uint(idPasien),
		IDJadwal:  uint(input.IDJadwal),
		Keluhan:   input.Keluhan,
		NoAntrian: lastAntrian + 1,
	}

	err := s.repo.SavePeriksaByPasien(&daftar)
	return &daftar, err
}

func (s *service) GetAntrianDokter(idDokter int, hariIni string) ([]domain.DaftarPoli, error) {
	return s.repo.FindPeriksaByDokter(idDokter, hariIni)
}

func (s *service) CreatePeriksa(input PeriksaInput) error {
	if len(input.Obat) == 0 {
		return errors.New("minimal harus meresepkan 1 jenis obat")
	}

	periksa := domain.Periksa{
		IDDaftarPoli: uint(input.IDDaftarPoli),
		TglPeriksa:   time.Now(),
		Catatan:      input.Catatan,
	}

	// Konversi ObatInput dari Handler menjadi DetailObat untuk Repository
	var obatRepoInputs []DetailObat
	for _, o := range input.Obat {
		obatRepoInputs = append(obatRepoInputs, DetailObat{
			IDObat: o.IDObat,
			Jumlah: o.Jumlah,
		})
	}

	return s.repo.SavePeriksaByDokterTx(&periksa, obatRepoInputs)
}
