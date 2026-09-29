package auth

import (
	"be-golang-poliklinik/internal/domain"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type MockRepository struct {
	MockFindByEmail func(email string) (*domain.User, error)
}
func (m *MockRepository) FindByEmail(email string) (*domain.User, error) { return m.MockFindByEmail(email) }
func (m *MockRepository) Save(user *domain.User) error { return nil }
func (m *MockRepository) NoRM(prefix string) (int64, error) { return 0, nil }
func (m *MockRepository) SaveBlacklist(token string, exp time.Time) error { return nil }
func (m *MockRepository) IsBlacklisted(token string) (bool, error) { return false, nil }

func TestLogin_SalahEmail(t *testing.T) {
	repo := &MockRepository{
		MockFindByEmail: func(email string) (*domain.User, error) {
			return nil, errors.New("not found") // Email tidak ditemukan
		},
	}
	service := NewService(repo)
	_, _, err := service.Login("salah@mail.com", "123456")
	if err == nil || err.Error() != "Email salah" {
		t.Errorf("Harusnya muncul error 'Email salah', tapi malah: %v", err)
	}
}

func TestLogin_SalahPassword(t *testing.T) {
	// Pura-puranya kita punya password "passwordbenar" yang sudah di-hash di DB
	hashed, _ := bcrypt.GenerateFromPassword([]byte("passwordbenar"), bcrypt.MinCost)
	repo := &MockRepository{
		MockFindByEmail: func(email string) (*domain.User, error) {
			return &domain.User{Email: email, Password: string(hashed)}, nil
		},
	}
	service := NewService(repo)
	_, _, err := service.Login("benar@mail.com", "passwordsalah")
	if err == nil || err.Error() != "Password salah" {
		t.Errorf("Harusnya muncul error 'Password salah', tapi malah: %v", err)
	}
}
