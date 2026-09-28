package auth

import (
	"be-golang-poliklinik/internal/domain"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Save(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	NoRM(mounthYear string) (int64, error)
	SaveBlacklist(token string, expiredAt time.Time) error
	IsBlacklisted(token string) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *repository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User

	err := r.db.Preload("Poli").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) NoRM(mounthYear string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Where("role = ? AND no_rm LIKE ?", domain.RolePasien, "RM-"+mounthYear+"%").Count(&count).Error
	return count, err
}

func (r *repository) SaveBlacklist(token string, expiredAt time.Time) error {
	blacklist := domain.BlacklistToken{
		Token:     token,
		ExpiredAt: expiredAt,
	}
	return r.db.Create(&blacklist).Error
}

func (r *repository) IsBlacklisted(token string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.BlacklistToken{}).Where("token = ?", token).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
