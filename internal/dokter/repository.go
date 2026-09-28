package dokter

import (
	"be-golang-poliklinik/internal/domain"

	"gorm.io/gorm"
)

type Repository interface {
	Save(dokter *domain.User) error
	FindAll() ([]domain.User, error)
	FindByID(id int) (*domain.User, error)
	Update(dokter *domain.User) error
	Delete(dokter *domain.User) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(dokter *domain.User) error {
	return r.db.Create(dokter).Error
}

func (r *repository) FindAll() ([]domain.User, error) {
	var dokters []domain.User
	err := r.db.Preload("Poli").Where("role = ?", domain.RoleDokter).Find(&dokters).Error
	return dokters, err
}

func (r *repository) FindByID(id int) (*domain.User, error) {
	var dokter domain.User
	err := r.db.Preload("Poli").Where("id = ? AND role = ?", id, domain.RoleDokter).First(&dokter).Error
	if err != nil {
		return nil, err
	}
	return &dokter, nil
}

func (r *repository) Update(dokter *domain.User) error {
	return r.db.Save(dokter).Error
}

func (r *repository) Delete(dokter *domain.User) error {
	return r.db.Delete(dokter).Error
}
