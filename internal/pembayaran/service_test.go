package pembayaran

import (
	"be-golang-poliklinik/internal/domain"
	"errors"
	"testing"
)

type MockRepository struct {
	MockGetPeriksaByID func(id int) (*domain.Periksa, error)
}
func (m *MockRepository) GetPeriksaByID(id int) (*domain.Periksa, error) { return m.MockGetPeriksaByID(id) }
func (m *MockRepository) GetTagihanAktif() ([]domain.Periksa, error) { return nil, nil }
func (m *MockRepository) SavePembayaran(pembayaran *domain.Pembayaran) error { return nil }
func (m *MockRepository) GetMenungguKonfirmasi() ([]domain.Pembayaran, error) { return nil, nil }
func (m *MockRepository) KonfirmasiPembayaran(id int) error { return nil }

func TestCreatePembayaran_GagalSudahLunas(t *testing.T) {
	repo := &MockRepository{
		MockGetPeriksaByID: func(id int) (*domain.Periksa, error) {
			// Simulasi DB bilang kalau datanya gak ada / udah lunas
			return nil, errors.New("tidak ditemukan atau sudah lunas") 
		},
	}
	service := NewService(repo)
	_, err := service.CreatePembayaran(PembayaranInput{IDPeriksa: 1, MetodeBayar: "Tunai"})
	
	if err == nil {
		t.Errorf("Harusnya sistem memblokir pembayaran ganda")
	}
	if err != nil && err.Error() != "tagihan tidak ditemukan atau sudah lunas" {
		t.Errorf("Diharapkan error tagihan lunas, tapi: %v", err)
	}
}
