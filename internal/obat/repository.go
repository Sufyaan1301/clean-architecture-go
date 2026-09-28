package obat

import (
	"be-golang-poliklinik/internal/domain"

	"gorm.io/gorm"
)

type Repository interface {
	Save(obat *domain.Obat) error
	FindAll() ([]domain.Obat, error)
	FindByID(id int) (*domain.Obat, error)
	Update(obat *domain.Obat) error
	Delete(obat *domain.Obat) error
	FindObatAvailable() ([]domain.Obat, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(obat *domain.Obat) error {
	return r.db.Create(obat).Error
}

func (r *repository) FindAll() ([]domain.Obat, error) {
	var obats []domain.Obat
	err := r.db.Find(&obats).Error
	return obats, err
}

func (r *repository) FindByID(id int) (*domain.Obat, error) {
	var obat domain.Obat
	err := r.db.Where("id = ?", id).First(&obat).Error
	if err != nil {
		return nil, err
	}
	return &obat, nil
}

func (r *repository) Update(obat *domain.Obat) error {
	return r.db.Save(obat).Error
}

func (r *repository) Delete(obat *domain.Obat) error {
	return r.db.Delete(obat).Error
}

func (r *repository) FindObatAvailable() ([]domain.Obat, error) {
	var obats []domain.Obat

	err := r.db.Where("stok > 0").Find(&obats).Error
	return obats, err
}
