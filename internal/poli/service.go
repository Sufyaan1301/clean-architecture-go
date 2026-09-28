package poli

import (
	"be-golang-poliklinik/internal/domain"
	"errors"
)

type Service interface {
	CreatePoli(input PoliInput) (*domain.Poli, error)
	GetAllPoli() ([]domain.Poli, error)
	GetPolibyID(id int) (*domain.Poli, error)
	UpdatePoli(id int, input PoliInput) (*domain.Poli, error)
	DeletePoli(id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

type PoliInput struct {
	NamaPoli   string `json:"nama_poli" binding:"required"`
	Keterangan string `json:"keterangan"`
}

func (s *service) CreatePoli(input PoliInput) (*domain.Poli, error) {
	poli := domain.Poli{
		NamaPoli:   input.NamaPoli,
		Keterangan: input.Keterangan,
	}

	err := s.repo.Save(&poli)
	return &poli, err
}

func (s *service) GetAllPoli() ([]domain.Poli, error) {
	return s.repo.FindAll()
}

func (s *service) GetPolibyID(id int) (*domain.Poli, error) {
	poli, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("data Poli tidak ditemukan")
	}
	return poli, nil
}

func (s *service) UpdatePoli(id int, input PoliInput) (*domain.Poli, error) {
	poli, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("data poli tidak ditemukan")
	}

	poli.NamaPoli = input.NamaPoli
	poli.Keterangan = input.Keterangan

	err = s.repo.Update(poli)
	return poli, err
}

func (s *service) DeletePoli(id int) error {
	poli, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("data poli tidak ditemukan")
	}

	return s.repo.Delete(poli)
}
