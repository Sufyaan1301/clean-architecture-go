package periksa

import (
	"be-golang-poliklinik/internal/domain"
	"errors"
	"testing"
)

// ==========================================
// 1. MOCKING (MEMBUAT DATABASE TIRUAN)
// ==========================================
type MockRepository struct {
	MockFindLastPeriksaByPasien func(idPasien int) error
	MockFindLastAntrian         func(idJadwal int) int
	MockSavePeriksaByPasien     func(daftar *domain.DaftarPoli) error
}

func (m *MockRepository) FindLastPeriksaByPasien(idPasien int) error {
	return m.MockFindLastPeriksaByPasien(idPasien)
}
func (m *MockRepository) FindLastAntrian(idJadwal int) int {
	return m.MockFindLastAntrian(idJadwal)
}
func (m *MockRepository) SavePeriksaByPasien(daftar *domain.DaftarPoli) error {
	return m.MockSavePeriksaByPasien(daftar)
}
// Fungsi di bawah ini dibiarkan kosong karena tidak dipakai di skenario tes ini
func (m *MockRepository) FindPeriksaByDokter(idDokter int, hariIni string) ([]domain.DaftarPoli, error) {
	return nil, nil
}
func (m *MockRepository) SavePeriksaByDokterTx(periksa *domain.Periksa, obatInputs []DetailObat) error {
	return nil
}

// ==========================================
// 2. SKENARIO TESTING
// ==========================================

// Tes 1: Pasien daftar antrian berhasil (Nomor antrian harus bertambah +1)
func TestCreateDaftarPoli_Berhasil(t *testing.T) {
	mockRepo := &MockRepository{
		MockFindLastPeriksaByPasien: func(idPasien int) error {
			return nil // Lolos dari blokir tagihan
		},
		MockFindLastAntrian: func(idJadwal int) int {
			return 10 // Pura-puranya antrian terakhir di DB adalah 10
		},
		MockSavePeriksaByPasien: func(daftar *domain.DaftarPoli) error {
			daftar.ID = 1 // Simulasi berhasil simpan ke DB
			return nil
		},
	}

	service := NewService(mockRepo)
	input := DaftarPoliInput{IDJadwal: 1, Keluhan: "Sakit Kepala"}

	hasil, err := service.CreateDaftarPoli(3, input)

	if err != nil {
		t.Errorf("Diharapkan berhasil, tapi mendapat error: %v", err)
	}

	// Cek apakah nomor antriannya 11 (10 + 1)
	if hasil.NoAntrian != 11 {
		t.Errorf("Diharapkan nomor antrian 11, tapi mendapat: %d", hasil.NoAntrian)
	}
}

// Tes 2: Pasien daftar antrian tapi DITOLAK karena punya tagihan menunggak
func TestCreateDaftarPoli_GagalAdaTagihan(t *testing.T) {
	mockRepo := &MockRepository{
		MockFindLastPeriksaByPasien: func(idPasien int) error {
			return errors.New("Pendaftaran ditolak: tagihan belum lunas") // DB Tiruan memblokir
		},
	}

	service := NewService(mockRepo)
	input := DaftarPoliInput{IDJadwal: 1, Keluhan: "Sakit Kepala"}

	hasil, err := service.CreateDaftarPoli(3, input)

	if err == nil {
		t.Errorf("Diharapkan gagal (error tagihan), tapi malah berhasil")
	}

	if hasil != nil {
		t.Errorf("Diharapkan hasil pendaftaran nil, tapi malah mendapat data")
	}
}

// Tes 3: Dokter submit periksa tapi lupa masukin obat (Validasi Array Kosong)
func TestCreatePeriksa_GagalTanpaObat(t *testing.T) {
	service := NewService(&MockRepository{}) // Repo tidak terpakai karena error divalidasi duluan

	input := PeriksaInput{
		IDDaftarPoli: 1,
		Catatan:      "Banyak istirahat",
		Obat:         []ObatInput{}, // Dokter lupa masukin obat (kosong)
	}

	err := service.CreatePeriksa(input)

	if err == nil {
		t.Errorf("Diharapkan ada error minimal 1 obat, tapi sistem meloloskannya")
	}

	expectedMsg := "minimal harus meresepkan 1 jenis obat"
	if err != nil && err.Error() != expectedMsg {
		t.Errorf("Diharapkan tulisan error '%s', tapi mendapat '%s'", expectedMsg, err.Error())
	}
}
