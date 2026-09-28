package dokter

import (
	"be-golang-poliklinik/internal/domain"
	"errors"
	"os"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateDokter(input DokterInput) (*domain.User, error)
	GetAllDokter() ([]domain.User, error)
	UpdateDokter(id int, input UpdateDokterInput) (*domain.User, error)
	GetDokterByID(id int) (*domain.User, error)
	DeleteDokter(id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

type DokterInput struct {
	Nama     string `form:"nama" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
	Alamat   string `form:"alamat" binding:"required"`
	NoKTP    string `form:"no_ktp" binding:"required"`
	NoHP     string `form:"no_hp" binding:"required"`
	IDPoli   uint   `form:"id_poli" binding:"required"`
	Foto     string
}

type UpdateDokterInput struct {
	Nama     string `form:"nama"`
	Email    string `form:"email"`
	Password string `form:"password"`
	Alamat   string `form:"alamat"`
	NoKTP    string `form:"no_ktp"`
	NoHP     string `form:"no_hp"`
	IDPoli   uint   `form:"id_poli"`
	Foto     string
}

func (s *service) CreateDokter(input DokterInput) (*domain.User, error) {
	if input.Password == "" {
		return nil, errors.New("password tidak boleh kosong")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return nil, errors.New("gagal mengenkripsi password")
	}

	dokter := domain.User{
		Nama:     input.Nama,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     domain.RoleDokter,
		Alamat:   input.Alamat,
		NoKTP:    input.NoKTP,
		NoHP:     input.NoHP,
		IDPoli:   &input.IDPoli,
		Foto:     input.Foto,
	}

	err = s.repo.Save(&dokter)
	return &dokter, err
}

func (s *service) GetAllDokter() ([]domain.User, error) {
	return s.repo.FindAll()
}

func (s *service) GetDokterByID(id int) (*domain.User, error) {
	dokter, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("data dokter tidak ditemukan")
	}
	return dokter, err
}

func (s *service) UpdateDokter(id int, input UpdateDokterInput) (*domain.User, error) {
	dokter, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("data dokter tidak ditemukan")
	}

	if input.Nama != "" {
		dokter.Nama = input.Nama
	}
	if input.Email != "" {
		dokter.Email = input.Email
	}
	if input.Alamat != "" {
		dokter.Alamat = input.Alamat
	}
	if input.NoKTP != "" {
		dokter.NoKTP = input.NoKTP
	}
	if input.NoHP != "" {
		dokter.NoHP = input.NoHP
	}
	if input.IDPoli != 0 {
		newPoliID := input.IDPoli
		dokter.IDPoli = &newPoliID
		dokter.Poli = nil
	}
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
		if err != nil {
			return nil, errors.New("gagal mengenkripsi password")
		}
		dokter.Password = string(hashedPassword)
	}

	if input.Foto != "" {
		if dokter.Foto != "" {
			os.Remove(dokter.Foto)
		}
		dokter.Foto = input.Foto
	}
	err = s.repo.Update(dokter)
	return dokter, err
}

func (s *service) DeleteDokter(id int) error {
	dokter, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("data dokter tidak ditemukan")
	}

	if dokter.Foto != "" {
		os.Remove(dokter.Foto)
	}

	return s.repo.Delete(dokter)
}
