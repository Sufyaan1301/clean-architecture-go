package obat

import (
	"be-golang-poliklinik/internal/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ObatHandler struct {
	service Service
}

func NewHandler(service Service) *ObatHandler {
	return &ObatHandler{service}
}

func (h *ObatHandler) CreateObat(c *gin.Context) {
	var input ObatInput
	if err := c.ShouldBind(&input); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
	}
	newObat, err := h.service.CreateObat(input)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	helper.Success(c, http.StatusCreated, "Obat Berhasil ditambahkan", newObat)
}

func (h *ObatHandler) GetAllObat(c *gin.Context) {
	polis, err := h.service.GetAllObat()
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "Data Obat berhasil diambil", polis)
}

func (h *ObatHandler) UpdatePoli(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "ID harus berupa angka")
		return
	}

	var input UpdateObatInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedPoli, err := h.service.UpdateObat(id, input)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "Data Obat berhasil diupdate!", updatedPoli)
}

func (h *ObatHandler) DeleteObat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "ID harus berupa angka")
		return
	}

	err = h.service.DeleteObat(id)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "Data Obat berhasil dihapus", nil)
}

func (h *ObatHandler) GetAllObatAvailable(c *gin.Context) {
	obatsavailable, err := h.service.GetAllObatAvailable()
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	formattedObatAvailable := FormatObatAvailableList(obatsavailable)

	helper.Success(c, http.StatusOK, "Data obat tersedia berhasil diambil", formattedObatAvailable)
}
