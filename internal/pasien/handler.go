package pasien

import (
	"be-golang-poliklinik/internal/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PasienHandler struct {
	service Service
}

func NewHandler(service Service) *PasienHandler {
	return &PasienHandler{service}
}

func (h *PasienHandler) GetAllPasien(c *gin.Context) {
	pasiens, err := h.service.GetAllPasien()
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	formattedPasien := FormatPasienList(pasiens)
	helper.Success(c, http.StatusOK, "Berhasil mengambil seluruh data pasien", formattedPasien)
}

func (h *PasienHandler) GetPasienByID(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "Parameter ID Pasien tidak valid")
		return
	}

	pasien, err := h.service.GetPasienByID(id)
	if err != nil {
		helper.Error(c, http.StatusNotFound, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "Berhasil mengambil detail data pasien", pasien)
}

func (h *PasienHandler) UpdatePasien(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "Parameter ID Pasien tidak valid")
		return
	}

	var input UpdatePasienInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.Error(c, http.StatusBadRequest, "Format input JSON tidak valid")
		return
	}

	updatedPasien, err := h.service.UpdatePasien(id, input)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "Berhasil memperbarui data pasien", FormatPasien(*updatedPasien))
}

func (h *PasienHandler) DeletePasien(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "Parameter ID Pasien tidak valid")
		return
	}

	err = h.service.DeletePasien(id)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "Berhasil menghapus data pasien beserta akunnya", nil)
}
