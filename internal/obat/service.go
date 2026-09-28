package obat

import (
	"be-golang-poliklinik/internal/domain"
	"errors"
)

type Service interface {
	CreateObat(input ObatInput) (*domain.Obat, error)
	GetAllObat() ([]domain.Obat, error)
	GetObatByID(id int) (*domain.Obat, error)
	UpdateObat(id int, input UpdateObatInput) (*domain.Obat, error)
	DeleteObat(id int) error
	GetAllObatAvailable() ([]domain.Obat, error)
}

type service struct {
	repo Repository
}

func (s *service) GetObatByID(id int) (*domain.Obat, error) {
	panic("unimplemented")
}

func NewService(repo Repository) *service {
	return &service{repo}
}

type ObatInput struct {
	NamaObat string `json:"nama_obat" binding:"required"`
	Kemasan  string `json:"kemasan" binding:"required"`
	Harga    int    `json:"harga" binding:"required,min=0"`
}

type UpdateObatInput struct {
	NamaObat string `json:"nama_obat"`
	Kemasan  string `json:"kemasan"`
	Harga    int    `json:"harga"`
	Stok     int    `json:"stok"`
}

func (s *service) CreateObat(input ObatInput) (*domain.Obat, error) {
	obat := domain.Obat{
		NamaObat: input.NamaObat,
		Kemasan:  input.Kemasan,
		Harga:    input.Harga,
	}
	err := s.repo.Save(&obat)
	return &obat, err
}

func (s *service) GetAllObat() ([]domain.Obat, error) {
	return s.repo.FindAll()
}

func (s *service) UpdateObat(id int, input UpdateObatInput) (*domain.Obat, error) {
	obat, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("data obat tidak ditemukan")
	}

	if input.NamaObat != "" {
		obat.NamaObat = input.NamaObat
	}
	if input.Kemasan != "" {
		obat.Kemasan = input.Kemasan
	}
	if input.Harga != 0 {
		obat.Harga = input.Harga
	}

	if input.Stok != 0 {
		obat.Stok = input.Stok
	}

	err = s.repo.Update(obat)
	return obat, err
}

func (s *service) DeleteObat(id int) error {
	obat, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("data obat tidak ditemukan")
	}

	return s.repo.Delete(obat)
}

func (s *service) GetAllObatAvailable() ([]domain.Obat, error) {
	return s.repo.FindObatAvailable()
}
