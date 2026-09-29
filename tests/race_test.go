package tests

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"be-golang-poliklinik/internal/config"
	"be-golang-poliklinik/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func TestRace_ConcurrentBooking(t *testing.T) {
	r := setupTestRouter()

	// 1. Setup Data Dummy Secara Cepat (Bypass API untuk setup demi kecepatan)
	poli := domain.Poli{NamaPoli: "Poli Balap", Keterangan: "Poli khusus test"}
	config.DB.Create(&poli)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.MinCost)
	dokter := domain.User{Nama: "Dr. Cepat", Email: "drcepat@mail.com", Password: string(hashedPassword), Role: domain.RoleDokter, IDPoli: &poli.ID}
	config.DB.Create(&dokter)

	timeMulai := time.Date(2026, 1, 1, 8, 0, 0, 0, time.Local)
	timeSelesai := time.Date(2026, 1, 1, 12, 0, 0, 0, time.Local)
	jadwal := domain.JadwalPeriksa{IDDokter: dokter.ID, Hari: "Selasa", JamMulai: timeMulai, JamSelesai: timeSelesai}
	config.DB.Create(&jadwal)

	// 2. Buat 10 Pasien & Ambil Tokennya
	var tokens []string
	for i := 1; i <= 10; i++ {
		email := "pasien" + strconv.Itoa(i) + "@mail.com"
		pasien := domain.User{
			Nama:     "Pasien " + strconv.Itoa(i),
			Email:    email,
			Password: string(hashedPassword),
			Role:     domain.RolePasien,
			NoRM:     "RM-" + strconv.Itoa(i),
			NoKTP:    "999888" + strconv.Itoa(i), // <-- Fix duplicate NoKTP
		}
		config.DB.Create(&pasien)

		// Hit Login API untuk dapat JWT
		w := sendJSON(r, "POST", "/api/auth/login", "", map[string]interface{}{"email": email, "password": "123"})
		tokens = append(tokens, parseJSONToken(w.Body.Bytes()))
	}

	// 3. FIRE CONCURRENT REQUESTS (RACE CONDITION SIMULATION)
	var wg sync.WaitGroup

	t.Log("Mulai Serangan! Menembakkan 10 request booking di milidetik yang SAMA PERSIS...")
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(token string, index int) {
			defer wg.Done()
			sendJSON(r, "POST", "/api/pasien/booking-periksa", token, map[string]interface{}{
				"id_jadwal": jadwal.ID,
				"keluhan":   "Sakit bersamaan dari goroutine " + strconv.Itoa(index),
			})
		}(tokens[i], i)
	}

	wg.Wait() // Tunggu semua goroutine (thread) selesai

	// 4. Validasi Hasil di Database
	var antrians []domain.DaftarPoli
	config.DB.Where("id_jadwal = ?", jadwal.ID).Order("no_antrian asc").Find(&antrians)

	t.Logf("Total antrian yang berhasil terdaftar: %d", len(antrians))

	if len(antrians) != 1 {
		t.Errorf("Keren! Sistem berhasil menahan serangan. Diharapkan hanya 1 antrian yang lolos, dan tersimpan %d", len(antrians))
	}

	// Cek Duplikat Nomor Antrian
	queueMap := make(map[int]bool)
	for _, a := range antrians {
		if queueMap[a.NoAntrian] {
			t.Fatalf("🚨 FATAL RACE CONDITION: Ditemukan duplikat nomor antrian: %d!", a.NoAntrian)
		}
		queueMap[a.NoAntrian] = true
		t.Logf("Sukses: Pasien ID %d mendapat Antrian No %d", a.IDPasien, a.NoAntrian)
	}
}
