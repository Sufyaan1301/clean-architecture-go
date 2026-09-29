package jadwalperiksa

import (
	"be-golang-poliklinik/internal/domain"
	"testing"
)

type MockRepository struct {
	MockFindByDokterAndHari func(idDokter int, hari string) ([]domain.JadwalPeriksa, error)
}
func (m *MockRepository) FindByDokterAndHari(id int, hari string) ([]domain.JadwalPeriksa, error) { return m.MockFindByDokterAndHari(id, hari) }
func (m *MockRepository) Save(j *domain.JadwalPeriksa) error { return nil }
func (m *MockRepository) FindByDokterID(id int) ([]domain.JadwalPeriksa, error) { return nil, nil }
func (m *MockRepository) FindByID(id int) (*domain.JadwalPeriksa, error) { return nil, nil }
func (m *MockRepository) Update(j *domain.JadwalPeriksa) error { return nil }
func (m *MockRepository) Delete(j *domain.JadwalPeriksa) error { return nil }
func (m *MockRepository) FindAll() ([]domain.JadwalPeriksa, error) { return nil, nil }

func TestCreateJadwal_GagalTabrakanWaktu(t *testing.T) {
	repo := &MockRepository{
		MockFindByDokterAndHari: func(id int, hari string) ([]domain.JadwalPeriksa, error) {
			jLama, _ := parseJam("08:00")
			jSelesai, _ := parseJam("12:00")
			return []domain.JadwalPeriksa{
				{Hari: "Senin", JamMulai: jLama, JamSelesai: jSelesai}, // Sudah ada jadwal 08.00 - 12.00
			}, nil
		},
	}
	service := NewService(repo)
	
	// Mencoba bikin jadwal baru jam 10.00 - 14.00 (bertabrakan di jam 10-12)
	input := JadwalPeriksaInput{Hari: "Senin", JamMulai: "10:00", JamSelesai: "14:00"}
	_, err := service.CreateJadwalPeriksa(2, input)
	
	if err == nil {
		t.Errorf("Harusnya sistem mendeteksi tabrakan waktu dan memblokirnya")
	}
	if err != nil && err.Error() != "Waktu bertabrakan dengan jadwal Anda yang lain di hari tersebut" {
		t.Errorf("Error messagenya tidak sesuai: %v", err)
	}
}
