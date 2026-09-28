package pasien

import (
	"be-golang-poliklinik/internal/domain"

	"gorm.io/gorm"
)

type Repository interface {
	FindAllPasien() ([]domain.User, error)
	FindByID(id int) (*domain.User, error)
	Update(pasien *domain.User) error
	Delete(pasien *domain.User) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAllPasien() ([]domain.User, error) {
	var pasiens []domain.User
	err := r.db.Where("role = ?", domain.RolePasien).Find(&pasiens).Error
	return pasiens, err
}

func (r *repository) FindByID(id int) (*domain.User, error) {
	var pasien domain.User
	err := r.db.Where("id = ? AND role = ?", id, domain.RolePasien).First(&pasien).Error
	if err != nil {
		return nil, err
	}
	return &pasien, nil
}

func (r *repository) Update(pasien *domain.User) error {
	return r.db.Save(pasien).Error
}

func (r *repository) Delete(pasien *domain.User) error {
	return r.db.Delete(pasien).Error
}
