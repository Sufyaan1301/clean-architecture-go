package pasien

import (
	"be-golang-poliklinik/internal/domain"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetAllPasien() ([]domain.User, error)
	GetPasienByID(id int) (*domain.User, error)
	UpdatePasien(id int, input UpdatePasienInput) (*domain.User, error)
	DeletePasien(id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

type UpdatePasienInput struct {
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Alamat   string `json:"alamat"`
	NoKTP    string `json:"no_ktp"`
	NoHP     string `json:"no_hp"`
}

func (s *service) GetAllPasien() ([]domain.User, error) {
	return s.repo.FindAllPasien()
}

func (s *service) GetPasienByID(id int) (*domain.User, error) {
	pasien, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("data pasien tidak ditemukan")
	}
	return pasien, nil
}

func (s *service) UpdatePasien(id int, input UpdatePasienInput) (*domain.User, error) {
	pasien, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("data pasien tidak ditemukan")
	}

	if input.Nama != "" {
		pasien.Nama = input.Nama
	}
	if input.Email != "" {
		pasien.Email = input.Email
	}
	if input.Alamat != "" {
		pasien.Alamat = input.Alamat
	}
	if input.NoKTP != "" {
		pasien.NoKTP = input.NoKTP
	}
	if input.NoHP != "" {
		pasien.NoHP = input.NoHP
	}
	
	if input.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
		pasien.Password = string(hashedPassword)
	}

	err = s.repo.Update(pasien)
	if err != nil {
		return nil, errors.New("gagal mengupdate data pasien (pastikan Email / NIK belum dipakai orang lain)")
	}

	return pasien, nil
}

func (s *service) DeletePasien(id int) error {
	pasien, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("data pasien tidak ditemukan")
	}

	err = s.repo.Delete(pasien)
	if err != nil {
		return errors.New("gagal menghapus data pasien")
	}

	return nil
}