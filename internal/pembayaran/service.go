package pembayaran

import (
	"be-golang-poliklinik/internal/domain"
	"errors"
	"time"
)

type Service interface {
	GetTagihanAktif() ([]domain.Periksa, error)
	CreatePembayaran(input PembayaranInput) (*domain.Pembayaran, error)
	GetMenungguKonfirmasi() ([]domain.Pembayaran, error)
	KonfirmasiPembayaran(idPembayaran int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

type PembayaranInput struct {
	IDPeriksa   int    `json:"id_periksa" binding:"required"`
	MetodeBayar string `json:"metode_bayar" binding:"required"` // Misal: Tunai, Transfer Bank, QRIS
}

func (s *service) GetTagihanAktif() ([]domain.Periksa, error) {
	return s.repo.GetTagihanAktif()
}

func (s *service) CreatePembayaran(input PembayaranInput) (*domain.Pembayaran, error) {
	periksa, err := s.repo.GetPeriksaByID(input.IDPeriksa)
	if err != nil {
		return nil, errors.New("tagihan tidak ditemukan atau sudah lunas")
	}
	pembayaran := domain.Pembayaran{
		IDPeriksa:   uint(input.IDPeriksa),
		TglBayar:    time.Now(),
		JumlahBayar: periksa.BiayaPeriksa,
		MetodeBayar: input.MetodeBayar,
	}

	err = s.repo.SavePembayaran(&pembayaran)
	if err != nil {
		return nil, err
	}

	return &pembayaran, nil
}

func (s *service) GetMenungguKonfirmasi() ([]domain.Pembayaran, error) {
	return s.repo.GetMenungguKonfirmasi()
}

func (s *service) KonfirmasiPembayaran(idPembayaran int) error {
	return s.repo.KonfirmasiPembayaran(idPembayaran)
}
