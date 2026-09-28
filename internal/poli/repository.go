package poli

import (
	"be-golang-poliklinik/internal/domain"

	"gorm.io/gorm"
)

type Repository interface {
	Save(poli *domain.Poli) error
	FindAll() ([]domain.Poli, error)
	FindByID(id int) (*domain.Poli, error)
	Update(poli *domain.Poli) error
	Delete(poli *domain.Poli) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(poli *domain.Poli) error {
	return r.db.Create(poli).Error
}

func (r *repository) FindAll() ([]domain.Poli, error) {
	var polis []domain.Poli
	err := r.db.Preload("Dokter").Find(&polis).Error
	return polis, err
}

func (r *repository) FindByID(id int) (*domain.Poli, error) {
	var poli domain.Poli
	err := r.db.Preload("Dokter").First(&poli, id).Error
	if err != nil {
		return nil, err
	}
	return &poli, nil
}

func (r *repository) Update(poli *domain.Poli) error {
	return r.db.Save(poli).Error
}

func (r *repository) Delete(poli *domain.Poli) error {
	return r.db.Delete(poli).Error
}