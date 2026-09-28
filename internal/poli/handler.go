package poli

import (
	"be-golang-poliklinik/internal/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PoliHandler struct {
	service Service
}

func NewHandler(service Service) *PoliHandler {
	return &PoliHandler{service}
}

func (h *PoliHandler) CreatePoli(c *gin.Context) {
	var input PoliInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	newPoli, err := h.service.CreatePoli(input)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusCreated, "Poli berhasil ditambahkan!", newPoli)
}

func (h *PoliHandler) GetAllPoli(c *gin.Context) {
	polis, err := h.service.GetAllPoli()
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	formattedPolis := FormatPoliList(polis)

	helper.Success(c, http.StatusOK, "Data poli berhasil diambil", formattedPolis)
}

func (h *PoliHandler) GetPoliByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "ID harus berupa angka")
		return
	}

	poli, err := h.service.GetPolibyID(id)
	if err != nil {
		helper.Error(c, http.StatusNotFound, err.Error())
		return
	}

	formatedDetailPolis := FormatPoliDetail(*poli)

	helper.Success(c, http.StatusOK, "Data poli berhasil diambil", formatedDetailPolis)
}

func (h *PoliHandler) UpdatePoli(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "ID harus berupa angka")
		return
	}

	var input PoliInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedPoli, err := h.service.UpdatePoli(id, input)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "Data poli berhasil diupdate!", updatedPoli)
}

func (h *PoliHandler) DeletePoli(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "ID harus berupa angka")
		return
	}

	err = h.service.DeletePoli(id)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "Data poli berhasil dihapus", nil)
}
