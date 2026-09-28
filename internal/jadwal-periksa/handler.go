package jadwalperiksa

import (
	"be-golang-poliklinik/internal/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JadwalPeriksaHandler struct {
	service Service
}

func NewHandler(service Service) *JadwalPeriksaHandler {
	return &JadwalPeriksaHandler{service}
}

func (h *JadwalPeriksaHandler) CreateJadwalPeriksa(c *gin.Context) {
	var input JadwalPeriksaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.Error(c, http.StatusBadRequest, "Format input tidak valid")
		return
	}

	tokenString := c.GetString("clean_token")
	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("rahasia_super_poliklinik"), nil
	})
	claims, _ := token.Claims.(jwt.MapClaims)
	idDokter := int(claims["user_id"].(float64))

	newJadwal, err := h.service.CreateJadwalPeriksa(idDokter, input)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusCreated, "Jadwal Periksa Berhasil ditambahkan", FormatJadwal(*newJadwal))
}

func (h *JadwalPeriksaHandler) GetJadwalPeriksa(c *gin.Context) {
	tokenString := c.GetString("clean_token")
	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("rahasia_super_poliklinik"), nil
	})
	claims, _ := token.Claims.(jwt.MapClaims)
	idDokter := int(claims["user_id"].(float64))

	jadwals, err := h.service.GetJadwalPeriksaByDokter(idDokter)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "Berhasil mengambil data Jadwal Periksa", FormatJadwalList(jadwals))
}

func (h *JadwalPeriksaHandler) UpdateJadwalPeriksa(c *gin.Context) {
	var input UpdateJadwalPeriksaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.Error(c, http.StatusBadRequest, "Format input tidak valid")
		return
	}

	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "ID Jadwal tidak valid")
		return
	}

	tokenString := c.GetString("clean_token")
	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("rahasia_super_poliklinik"), nil
	})
	claims, _ := token.Claims.(jwt.MapClaims)
	idDokter := int(claims["user_id"].(float64))

	updatedJadwal, err := h.service.UpdateJadwalPeriksa(id, idDokter, input)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "Jadwal Periksa Berhasil diubah", FormatJadwal(*updatedJadwal))
}

func (h *JadwalPeriksaHandler) DeleteJadwalPeriksa(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "ID Jadwal tidak valid")
		return
	}

	tokenString := c.GetString("clean_token")
	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("rahasia_super_poliklinik"), nil
	})
	claims, _ := token.Claims.(jwt.MapClaims)
	idDokter := int(claims["user_id"].(float64))

	err = h.service.DeleteJadwalPeriksa(id, idDokter)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "Berhasil menghapus jadwal periksa", nil)
}

func (h *JadwalPeriksaHandler) GetAllJadwalPeriksa(c *gin.Context) {
	jadwals, err := h.service.GetAllJadwalPeriksa()
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Gagal mengambil daftar jadwal praktek")
		return
	}

	helper.Success(c, http.StatusOK, "Berhasil mengambil seluruh daftar jadwal praktek dokter", FormatJadwalList(jadwals))
}


