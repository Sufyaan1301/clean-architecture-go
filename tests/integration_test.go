package tests

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"


	"be-golang-poliklinik/internal/auth"
	"be-golang-poliklinik/internal/config"
	"be-golang-poliklinik/internal/dokter"
	jadwalperiksa "be-golang-poliklinik/internal/jadwal-periksa"
	"be-golang-poliklinik/internal/obat"
	"be-golang-poliklinik/internal/pasien"
	"be-golang-poliklinik/internal/pembayaran"
	"be-golang-poliklinik/internal/periksa"
	"be-golang-poliklinik/internal/poli"
	"be-golang-poliklinik/internal/riwayat"
	"be-golang-poliklinik/internal/router"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestRouter() *gin.Engine {
	dsnInit := "root:@tcp(127.0.0.1:3306)/poliklinik_golang_db?charset=utf8mb4&parseTime=True&loc=Local"
	dbInit, _ := gorm.Open(mysql.Open(dsnInit), &gorm.Config{})
	dbInit.Exec("CREATE DATABASE IF NOT EXISTS poliklinik_golang_db_test")

	dsnTest := "root:@tcp(127.0.0.1:3306)/poliklinik_golang_db_test?charset=utf8mb4&parseTime=True&loc=Local"
	dbTest, _ := gorm.Open(mysql.Open(dsnTest), &gorm.Config{})

	config.DB = dbTest
	config.MigrateFresh()
	config.SeedAdmin() // Buat Akun Admin (admin@poliklinik.com / admin123)

	authRepo := auth.NewRepository(config.DB)
	authService := auth.NewService(authRepo)

	h := &router.Handlers{
		Auth:          auth.NewHandler(authService),
		Dokter:        dokter.NewHandler(dokter.NewService(dokter.NewRepository(config.DB))),
		Poli:          poli.NewHandler(poli.NewService(poli.NewRepository(config.DB))),
		Obat:          obat.NewHandler(obat.NewService(obat.NewRepository(config.DB))),
		Pasien:        pasien.NewHandler(pasien.NewService(pasien.NewRepository(config.DB))),
		JadwalPeriksa: jadwalperiksa.NewHandler(jadwalperiksa.NewService(jadwalperiksa.NewRepository(config.DB))),
		Periksa:       periksa.NewHandler(periksa.NewService(periksa.NewRepository(config.DB))),
		Riwayat:       riwayat.NewHandler(riwayat.NewService(riwayat.NewRepository(config.DB))),
		Pembayaran:    pembayaran.NewHandler(pembayaran.NewService(pembayaran.NewRepository(config.DB))),
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	router.SetupRoutes(r, h, authService)

	return r
}

func parseJSONToken(body []byte) string {
	var res map[string]interface{}
	json.Unmarshal(body, &res)
	data := res["data"].(map[string]interface{})
	return data["access_token"].(string)
}

func sendJSON(r *gin.Engine, method, path, token string, payload map[string]interface{}) *httptest.ResponseRecorder {
	var req *http.Request
	if payload != nil {
		body, _ := json.Marshal(payload)
		req, _ = http.NewRequest(method, path, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func sendMultipart(r *gin.Engine, method, path, token string, fields map[string]string) *httptest.ResponseRecorder {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, val := range fields {
		writer.WriteField(key, val)
	}

	// Buat file palsu untuk lolos validasi foto
	part, _ := writer.CreateFormFile("foto", "dummy.jpg")
	part.Write([]byte("fake image content"))

	writer.Close()

	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func parseID(body []byte) int {
	var res map[string]interface{}
	json.Unmarshal(body, &res)
	data := res["data"].(map[string]interface{})
	return int(data["id"].(float64))
}

func parseIDPembayaran(body []byte) int {
	var res map[string]interface{}
	json.Unmarshal(body, &res)
	data := res["data"].(map[string]interface{})
	return int(data["id_pembayaran"].(float64)) // Di formatter tadi kita menamainya id_pembayaran
}

func TestE2E_GoldenFlow(t *testing.T) {
	r := setupTestRouter()

	// 1. ADMIN LOGIN
	w1 := sendJSON(r, "POST", "/api/auth/login", "", map[string]interface{}{
		"email":    "admin@poliklinik.com",
		"password": "admin123",
	})
	if w1.Code != 200 { t.Fatalf("1. Admin Login Gagal: %v", w1.Body.String()) }
	adminToken := parseJSONToken(w1.Body.Bytes())

	// 2. ADMIN CREATE POLI
	w2 := sendJSON(r, "POST", "/api/admin/poli", adminToken, map[string]interface{}{
		"nama_poli": "Poli Gigi",
		"keterangan": "Poli Kesehatan Gigi",
	})
	if w2.Code != 201 { t.Fatalf("2. Create Poli Gagal: %v", w2.Body.String()) }

	// 3. ADMIN CREATE DOKTER (Pake Multipart Form)
	w3 := sendMultipart(r, "POST", "/api/admin/dokter", adminToken, map[string]string{
		"nama":     "Dr. Budi",
		"email":    "drbudi@mail.com",
		"password": "password123",
		"alamat":   "Jl. Sehat No. 1",
		"no_ktp":   "1111111111111111",
		"no_hp":    "081111111111",
		"id_poli":  "1",
	})
	if w3.Code != 201 { t.Fatalf("3. Create Dokter Gagal: %v", w3.Body.String()) }

	// 4. ADMIN CREATE OBAT & RESTOCK
	w4 := sendJSON(r, "POST", "/api/admin/obat", adminToken, map[string]interface{}{
		"nama_obat": "Paracetamol",
		"kemasan":   "Strip",
		"harga":     15000,
	})
	if w4.Code != 201 { t.Fatalf("4. Create Obat Gagal: %v", w4.Body.String()) }
	idObat := parseID(w4.Body.Bytes())

	sendJSON(r, "PUT", "/api/admin/obat/1", adminToken, map[string]interface{}{"stok": 50})

	// 5. DOKTER LOGIN
	w5 := sendJSON(r, "POST", "/api/auth/login", "", map[string]interface{}{
		"email":    "drbudi@mail.com",
		"password": "password123",
	})
	if w5.Code != 200 { t.Fatalf("5. Dokter Login Gagal: %v", w5.Body.String()) }
	dokterToken := parseJSONToken(w5.Body.Bytes())

	// 🛑 [SKENARIO NEGATIF 1] Dokter Ceroboh: Jam selesai lebih dulu dari jam mulai
	wNeg1 := sendJSON(r, "POST", "/api/dokter/jadwal-periksa", dokterToken, map[string]interface{}{
		"hari": "Selasa", "jam_mulai": "12:00", "jam_selesai": "08:00",
	})
	if wNeg1.Code != 400 { t.Fatalf("Sistem Kebobolan! Jadwal mundur bisa lolos: %v", wNeg1.Body.String()) }
	t.Log("✔️ [Tahan Serangan] Jadwal waktu mundur berhasil ditolak!")

	// 6. DOKTER CREATE JADWAL (Benar)
	w6 := sendJSON(r, "POST", "/api/dokter/jadwal-periksa", dokterToken, map[string]interface{}{
		"hari":        "Senin",
		"jam_mulai":   "08:00",
		"jam_selesai": "12:00",
	})
	if w6.Code != 201 { t.Fatalf("6. Create Jadwal Gagal: %v", w6.Body.String()) }
	idJadwal := parseID(w6.Body.Bytes())

	// 7 & 8. PASIEN REGISTER & LOGIN
	sendJSON(r, "POST", "/api/auth/register", "", map[string]interface{}{
		"nama": "Pasien Nakal", "email": "nakal@mail.com", "password": "password123", "no_ktp": "123",
	})
	w8 := sendJSON(r, "POST", "/api/auth/login", "", map[string]interface{}{
		"email": "nakal@mail.com", "password": "password123",
	})
	pasienToken := parseJSONToken(w8.Body.Bytes())

	// 9. PASIEN BOOKING (DAFTAR POLI) Pertama
	w9 := sendJSON(r, "POST", "/api/pasien/booking-periksa", pasienToken, map[string]interface{}{
		"id_jadwal": idJadwal, "keluhan": "Sakit Kepala",
	})
	if w9.Code != 201 { t.Fatalf("9. Booking Gagal: %v", w9.Body.String()) }
	idDaftarPoli := parseID(w9.Body.Bytes())

	// 🛑 [SKENARIO NEGATIF 2] Pasien Maruk: Coba daftar lagi padahal belum diperiksa
	wNeg2 := sendJSON(r, "POST", "/api/pasien/booking-periksa", pasienToken, map[string]interface{}{
		"id_jadwal": idJadwal, "keluhan": "Sakit Perut Juga",
	})
	if wNeg2.Code != 400 { t.Fatalf("Sistem Kebobolan! Pasien bisa dobel antrian: %v", wNeg2.Body.String()) }
	t.Log("✔️ [Tahan Serangan] Pendaftaran dobel berhasil ditolak!")

	// 🛑 [SKENARIO NEGATIF 3] Dokter Halu: Resep 1000 obat padahal stok cuma 50
	wNeg3 := sendJSON(r, "POST", "/api/dokter/periksa", dokterToken, map[string]interface{}{
		"id_daftar_poli": idDaftarPoli, "catatan": "Kasih 1000 obat",
		"obat": []map[string]interface{}{{"id_obat": idObat, "jumlah": 1000}},
	})
	if wNeg3.Code != 400 { t.Fatalf("Sistem Kebobolan! Dokter meresepkan obat melebihi stok: %v", wNeg3.Body.String()) }
	t.Log("✔️ [Tahan Serangan] Resep melebihi kapasitas gudang ditolak!")

	// 10. DOKTER SUBMIT PERIKSA (Benar - Minta 2 obat)
	w10 := sendJSON(r, "POST", "/api/dokter/periksa", dokterToken, map[string]interface{}{
		"id_daftar_poli": idDaftarPoli, "catatan": "Sakit biasa",
		"obat": []map[string]interface{}{{"id_obat": idObat, "jumlah": 2}},
	})
	if w10.Code != 201 { t.Fatalf("10. Submit Periksa Gagal: %v", w10.Body.String()) }
	
	// GET ID PERIKSA DARI RIWAYAT PASIEN (Karena POST Periksa mengembalikan nil)
	w10b := sendJSON(r, "GET", "/api/pasien/riwayat", pasienToken, nil)
	var riwayatRes map[string]interface{}
	json.Unmarshal(w10b.Body.Bytes(), &riwayatRes)
	riwayatData := riwayatRes["data"].([]interface{})
	firstRiwayat := riwayatData[0].(map[string]interface{})
	idPeriksa := int(firstRiwayat["id"].(float64))

	// 11. PASIEN BAYAR DI KASIR (Tapi belum dikonfirmasi admin!)
	w11 := sendJSON(r, "POST", "/api/pasien/pembayaran", pasienToken, map[string]interface{}{
		"id_periksa": idPeriksa, "metode_bayar": "Tunai",
	})
	if w11.Code != 201 { t.Fatalf("11. Pembayaran Gagal: %v", w11.Body.String()) }

	// 🛑 [SKENARIO NEGATIF 4] Pasien Tukang Ngutang: Status 'Menunggu Konfirmasi', coba daftar antrian baru
	wNeg4 := sendJSON(r, "POST", "/api/pasien/booking-periksa", pasienToken, map[string]interface{}{
		"id_jadwal": idJadwal, "keluhan": "Masih sakit kepala",
	})
	if wNeg4.Code != 400 { t.Fatalf("Sistem Kebobolan! Pasien ngutang bisa daftar lagi: %v", wNeg4.Body.String()) }
	t.Log("✔️ [Tahan Serangan] Pasien ngutang (belum Lunas) dilarang berobat lagi!")

	// 12. ADMIN KONFIRMASI LUNAS
	w12 := sendJSON(r, "PUT", "/api/admin/pembayaran/konfirmasi/1", adminToken, nil) 
	if w12.Code != 200 { t.Fatalf("12. Konfirmasi Gagal: %v", w12.Body.String()) }

	// 🛑 [SKENARIO NEGATIF 5] Hacker: Pasien nakal mau bikin Poli via API Admin
	wNeg5 := sendJSON(r, "POST", "/api/admin/poli", pasienToken, map[string]interface{}{
		"nama_poli": "Poli Hacker",
	})
	if wNeg5.Code == 201 || wNeg5.Code == 200 { t.Fatalf("FATAL! Pasien tembus jadi Admin!") }
	t.Log("✔️ [Tahan Serangan] Akses ilegal role ditahan oleh Middleware keamanan!")

	t.Logf("✅✅✅ GOLDEN FLOW + 5 SERANGAN NEGATIF BERHASIL DITAHAN SISTEM!")
}
