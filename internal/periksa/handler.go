package periksa

import (
	"be-golang-poliklinik/internal/helper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type PeriksaHandler struct {
	service Service
}

func NewHandler(service Service) *PeriksaHandler {
	return &PeriksaHandler{service}
}

func getUserID(c *gin.Context) int {
	tokenString := c.GetString("clean_token")
	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("rahasia_super_poliklinik"), nil
	})
	claims, _ := token.Claims.(jwt.MapClaims)
	return int(claims["user_id"].(float64))
}

func (h *PeriksaHandler) CreateDaftarPoli(c *gin.Context) {
	var input DaftarPoliInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.Error(c, http.StatusBadRequest, "Data pendaftaran tidak lengkap")
		return
	}

	idPasien := getUserID(c)

	newDaftar, err := h.service.CreateDaftarPoli(idPasien, input)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	formattedDaftar := FormatAntrian(*newDaftar)
	helper.Success(c, http.StatusCreated, "Berhasil mendaftar antrian", formattedDaftar)
}

func (h *PeriksaHandler) GetAntrianDokter(c *gin.Context) {
	idDokter := getUserID(c)

	days := map[string]string{
		"Sunday": "Minggu", "Monday": "Senin", "Tuesday": "Selasa",
		"Wednesday": "Rabu", "Thursday": "Kamis", "Friday": "Jumat", "Saturday": "Sabtu",
	}
	hariIni := days[time.Now().Weekday().String()]

	antrian, err := h.service.GetAntrianDokter(idDokter, hariIni)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Gagal mengambil daftar antrian")
		return
	}

	// Bungkus dengan formatter agar "Jadwal" hilang dan JSON lebih bersih
	formattedAntrian := FormatAntrianList(antrian)

	helper.Success(c, http.StatusOK, "Berhasil memuat daftar pasien hari ini", formattedAntrian)
}

func (h *PeriksaHandler) CreatePeriksa(c *gin.Context) {
	var input PeriksaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.Error(c, http.StatusBadRequest, "Format pemeriksaan tidak valid")
		return
	}

	err := h.service.CreatePeriksa(input)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusCreated, "Pemeriksaan selesai disimpan", nil)
}
