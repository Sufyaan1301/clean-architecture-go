package dokter

import (
	"be-golang-poliklinik/internal/helper"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type DokterHandler struct {
	service Service
}

func NewHandler(service Service) *DokterHandler {
	return &DokterHandler{service}
}

func (h *DokterHandler) CreateDokter(c *gin.Context) {
	var input DokterInput
	if err := c.ShouldBind(&input); err != nil {
		helper.Error(c, http.StatusBadRequest, helper.FormatValidatorError(err))
		return
	}

	file, errFile := c.FormFile("foto")
	if errFile != nil {
		helper.Error(c, http.StatusBadRequest, "Foto dokter wajib di-upload!")
		return
	}

	if file.Size > 2*1024*1024 {
		helper.Error(c, http.StatusBadRequest, "Ukuran gambar tidak boleh lebih dari 2 MB")
		return
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		helper.Error(c, http.StatusBadRequest, "Format file tidak diizinkan. Hanya boleh JPG, JPEG, atau PNG")
		return
	}
	os.MkdirAll("./uploads", os.ModePerm)
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := "uploads/" + newFileName

	if errSave := c.SaveUploadedFile(file, filePath); errSave != nil {
		helper.Error(c, http.StatusInternalServerError, "Gagal menyimpan gambar baru")
		return
	}

	input.Foto = filePath

	newDokter, err := h.service.CreateDokter(input)
	if err != nil {
		if input.Foto != "" {
			os.Remove(input.Foto)
		}
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusCreated, "Dokter berhasil ditambahkan!", newDokter)
}

func (h *DokterHandler) GetAllDokter(c *gin.Context) {
	dokters, err := h.service.GetAllDokter()
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	formattedDokter := FormatDokterList(dokters)

	helper.Success(c, http.StatusOK, "Data dokter berhasil diambil", formattedDokter)
}

func (h *DokterHandler) GetDokterByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "ID harus berupa angka")
		return
	}
	dokters, err := h.service.GetDokterByID(id)
	if err != nil {
		helper.Error(c, http.StatusNotFound, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "Data dokter berhasil diambil", dokters)
}

func (h *DokterHandler) UpdateDokter(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "ID harus berupa angka")
		return
	}

	var input UpdateDokterInput
	if err := c.ShouldBind(&input); err != nil {
		helper.Error(c, http.StatusBadRequest, helper.FormatValidatorError(err))
		return
	}
	file, errFile := c.FormFile("foto")
	if errFile == nil {
		if file.Size > 2*1024*1024 {
			helper.Error(c, http.StatusBadRequest, "Ukuran gambar tidak boleh lebih dari 2 MB")
			return
		}
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			helper.Error(c, http.StatusBadRequest, "Format file tidak diizinkan. Hanya boleh JPG, JPEG, atau PNG")
			return
		}
		os.MkdirAll("./uploads", os.ModePerm)
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		filePath := "uploads/" + newFileName

		if errSave := c.SaveUploadedFile(file, filePath); errSave != nil {
			helper.Error(c, http.StatusInternalServerError, "Gagal menyimpan gambar baru")
			return
		}

		input.Foto = filePath
	}
	updatedDokter, err := h.service.UpdateDokter(id, input)
	if err != nil {
		if input.Foto != "" {
			os.Remove(input.Foto)
		}
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "Data dokter berhasil diupdate!", updatedDokter)
}

func (h *DokterHandler) DeleteDokter(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "ID harus berupa angka")
		return
	}

	err = h.service.DeleteDokter(id)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "Data dokter berhasil dihapus", nil)
}
