package riwayat

import "be-golang-poliklinik/internal/domain"

type Service interface {
	GetRiwayatDokter(idDokter int) ([]domain.Periksa, error)
	GetRiwayatDokterByID(idPeriksa int, idDokter int) (*domain.Periksa, error)
	GetRiwayatPasien(idPasien int) ([]domain.Periksa, error)
	GetRiwayatPasienByID(idPeriksa int, idPasien int) (*domain.Periksa, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

func (s *service) GetRiwayatDokter(idDokter int) ([]domain.Periksa, error) {
	return s.repo.GetRiwayatDokter(idDokter)
}

func (s *service) GetRiwayatPasien(idPasien int) ([]domain.Periksa, error) {
	return s.repo.GetRiwayatPasien(idPasien)
}

func (s *service) GetRiwayatDokterByID(idPeriksa int, idDokter int) (*domain.Periksa, error) {
	return s.repo.GetRiwayatDokterByID(idPeriksa, idDokter)
}

func (s *service) GetRiwayatPasienByID(idPeriksa int, idPasien int) (*domain.Periksa, error) {
	return s.repo.GetRiwayatPasienByID(idPeriksa, idPasien)
}