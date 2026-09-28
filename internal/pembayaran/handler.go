package pembayaran

import (
	"be-golang-poliklinik/internal/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PembayaranHandler struct {
	service Service
}

func NewHandler(service Service) *PembayaranHandler {
	return &PembayaranHandler{service}
}

func (h *PembayaranHandler) GetTagihanAktif(c *gin.Context) {
	tagihans, err := h.service.GetTagihanAktif()
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Gagal mengambil data tagihan aktif")
		return
	}

	formatted := FormatTagihanList(tagihans)
	helper.Success(c, http.StatusOK, "Berhasil memuat daftar tagihan kasir", formatted)
}

func (h *PembayaranHandler) CreatePembayaran(c *gin.Context) {
	var input PembayaranInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.Error(c, http.StatusBadRequest, "Format pembayaran tidak valid")
		return
	}

	pembayaran, err := h.service.CreatePembayaran(input)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	formatted := FormatPembayaran(*pembayaran)
	helper.Success(c, http.StatusCreated, "Berhasil mengirim permintaan pembayaran", formatted)
}

func (h *PembayaranHandler) GetMenungguKonfirmasi(c *gin.Context) {
	pembayarans, err := h.service.GetMenungguKonfirmasi()
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Gagal mengambil data")
		return
	}

	formatted := FormatPembayaranList(pembayarans)
	helper.Success(c, http.StatusOK, "Berhasil memuat daftar konfirmasi", formatted)
}

func (h *PembayaranHandler) KonfirmasiPembayaran(c *gin.Context) {
	idParam := c.Param("id")
	idPembayaran, err := strconv.Atoi(idParam)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "ID Pembayaran harus berupa angka")
		return
	}

	err = h.service.KonfirmasiPembayaran(idPembayaran)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "Gagal mengonfirmasi pembayaran")
		return
	}

	helper.Success(c, http.StatusOK, "Pembayaran berhasil dikonfirmasi (Lunas)", nil)
}
