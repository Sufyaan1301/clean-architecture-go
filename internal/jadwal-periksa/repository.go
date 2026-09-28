package jadwalperiksa

import (
	"be-golang-poliklinik/internal/domain"

	"gorm.io/gorm"
)

type Repository interface {
	Save(jadwal *domain.JadwalPeriksa) error
	FindByDokterID(idDokter int) ([]domain.JadwalPeriksa, error)
	FindByID(id int) (*domain.JadwalPeriksa, error)
	Update(jadwal *domain.JadwalPeriksa) error
	Delete(jadwal *domain.JadwalPeriksa) error
	FindByDokterAndHari(idDokter int, hari string) ([]domain.JadwalPeriksa, error)
	FindAll() ([]domain.JadwalPeriksa, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(jadwal *domain.JadwalPeriksa) error {
	err := r.db.Create(jadwal).Error
	if err == nil {
		r.db.Preload("Dokter").First(jadwal, jadwal.ID)
	}

	return err
}

func (r *repository) FindByDokterID(idDokter int) ([]domain.JadwalPeriksa, error) {
	var jadwals []domain.JadwalPeriksa
	err := r.db.Where("id_dokter = ?", idDokter).Find(&jadwals).Error
	return jadwals, err
}

func (r *repository) FindByID(id int) (*domain.JadwalPeriksa, error) {
	var jadwal domain.JadwalPeriksa
	err := r.db.Where("id = ?", id).First(&jadwal).Error
	if err != nil {
		return nil, err
	}
	return &jadwal, nil
}

func (r *repository) Update(jadwal *domain.JadwalPeriksa) error {
	return r.db.Save(jadwal).Error
}

func (r *repository) Delete(jadwal *domain.JadwalPeriksa) error {
	return r.db.Delete(jadwal).Error
}

func (r *repository) FindByDokterAndHari(idDokter int, hari string) ([]domain.JadwalPeriksa, error) {
	var jadwals []domain.JadwalPeriksa
	err := r.db.Where("id_dokter = ? AND hari = ?", idDokter, hari).Find(&jadwals).Error
	return jadwals, err
}

func (r *repository) FindAll() ([]domain.JadwalPeriksa, error) {
	var jadwals []domain.JadwalPeriksa
	err := r.db.
		Preload("Dokter").
		Preload("Dokter.Poli").
		Find(&jadwals).Error
	
	return jadwals, err
}
