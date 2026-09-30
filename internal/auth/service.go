package auth

import (
	"be-golang-poliklinik/internal/domain"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(input PasienRegisterInput) (*domain.User, error)
	Login(email, password string) (string, *domain.User, error)
	Logout(tokenString string) error
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return []byte("rahasia_super_poliklinik") // Fallback
	}
	return []byte(secret)
}

type PasienRegisterInput struct {
	Nama     string `json:"nama" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Alamat   string `json:"alamat"`
	NoKTP    string `json:"no_ktp" binding:"required"`
	NoHP     string `json:"no_hp"`
}

func (s *service) Register(input PasienRegisterInput) (*domain.User, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	currentTime := time.Now()
	monthYearStr := currentTime.Format("200601")
	count, _ := s.repo.NoRM(monthYearStr)
	generateNoRM := fmt.Sprintf("RM-%s-%04d", monthYearStr, count+1)

	user := domain.User{
		Nama:     input.Nama,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     domain.RolePasien,
		Alamat:   input.Alamat,
		NoKTP:    input.NoKTP,
		NoHP:     input.NoHP,
		NoRM:     generateNoRM,
	}

	err := s.repo.Save(&user)
	return &user, err
}

func (s *service) Login(email, password string) (string, *domain.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("Email salah")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, errors.New("Password salah")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"nama":    user.Nama,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	})

	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		return "", nil, err
	}
	return tokenString, user, nil

}

func (s *service) ValidateToken(tokenString string) (*jwt.Token, error) {
	isBlacklisted, err := s.repo.IsBlacklisted(tokenString)
	if err != nil || isBlacklisted {
		return nil, errors.New("token sudah tidak berlaku, silakan login kembali")
	}
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})
}

func (s *service) Logout(tokenString string) error {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return errors.New("token tidak valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("klaim token tidak terbaca")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("token tidak memiliki jangka waktu ekspirasi")
	}
	expiredAt := time.Unix(int64(expFloat), 0)

	return s.repo.SaveBlacklist(tokenString, expiredAt)
}
