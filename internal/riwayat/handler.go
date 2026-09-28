package riwayat

import (
	"be-golang-poliklinik/internal/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type RiwayatHandler struct {
	service Service
}

func NewHandler(service Service) *RiwayatHandler {
	return &RiwayatHandler{service}
}

// Helper baca token
func getUserID(c *gin.Context) int {
	tokenString := c.GetString("clean_token")
	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("rahasia_super_poliklinik"), nil
	})
	claims, _ := token.Claims.(jwt.MapClaims)
	return int(claims["user_id"].(float64))
}

// ============================================
// HANDLER LIST (DAFTAR)
// ============================================

func (h *RiwayatHandler) GetRiwayatDokter(c *gin.Context) {
	idDokter := getUserID(c)

	riwayats, err := h.service.GetRiwayatDokter(idDokter)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Gagal memuat riwayat dokter")
		return
	}

	formatted := FormatRiwayatDokterList(riwayats)
	helper.Success(c, http.StatusOK, "Berhasil memuat riwayat rekam medis pasien", formatted)
}

func (h *RiwayatHandler) GetRiwayatPasien(c *gin.Context) {
	idPasien := getUserID(c)

	riwayats, err := h.service.GetRiwayatPasien(idPasien)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Gagal memuat riwayat pasien")
		return
	}

	formatted := FormatRiwayatPasienList(riwayats)
	helper.Success(c, http.StatusOK, "Berhasil memuat riwayat kunjungan", formatted)
}

// ============================================
// HANDLER DETAIL (1 DATA)
// ============================================

func (h *RiwayatHandler) GetRiwayatDetailDokter(c *gin.Context) {
	idDokter := getUserID(c)
	idPeriksa, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	riwayat, err := h.service.GetRiwayatDokterByID(idPeriksa, idDokter)
	if err != nil {
		helper.Error(c, http.StatusNotFound, "Data detail riwayat tidak ditemukan")
		return
	}

	formatted := FormatRiwayatDokter(*riwayat)
	helper.Success(c, http.StatusOK, "Berhasil memuat detail riwayat", formatted)
}

func (h *RiwayatHandler) GetRiwayatDetailPasien(c *gin.Context) {
	idPasien := getUserID(c)
	idPeriksa, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	riwayat, err := h.service.GetRiwayatPasienByID(idPeriksa, idPasien)
	if err != nil {
		helper.Error(c, http.StatusNotFound, "Data detail riwayat tidak ditemukan")
		return
	}

	formatted := FormatRiwayatPasien(*riwayat)
	helper.Success(c, http.StatusOK, "Berhasil memuat detail riwayat", formatted)
}
