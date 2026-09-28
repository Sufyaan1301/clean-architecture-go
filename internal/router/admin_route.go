package router

import (
	"be-golang-poliklinik/internal/auth"
	"be-golang-poliklinik/internal/middleware"

	"github.com/gin-gonic/gin"
)

func setupAdminRoutes(api *gin.RouterGroup, h *Handlers, authService auth.Service) {
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(authService))
	admin.Use(middleware.RoleBlockMiddleware("admin"))
	{
		admin.POST("/dokter", h.Dokter.CreateDokter)
		admin.GET("/dokter", h.Dokter.GetAllDokter)
		admin.GET("/dokter/:id", h.Dokter.GetDokterByID)
		admin.PUT("/dokter/:id", h.Dokter.UpdateDokter)
		admin.DELETE("/dokter/:id", h.Dokter.DeleteDokter)

		admin.POST("/poli", h.Poli.CreatePoli)
		admin.GET("/poli", h.Poli.GetAllPoli)
		admin.GET("/poli/:id", h.Poli.GetPoliByID)
		admin.PUT("/poli/:id", h.Poli.UpdatePoli)
		admin.DELETE("/poli/:id", h.Poli.DeletePoli)

		admin.GET("/obat", h.Obat.GetAllObat)
		admin.POST("/obat", h.Obat.CreateObat)
		admin.PUT("/obat/:id", h.Obat.UpdatePoli)
		admin.DELETE("/obat/:id", h.Obat.DeleteObat)

		admin.GET("/pasien", h.Pasien.GetAllPasien)
		admin.GET("/pasien/:id", h.Pasien.GetPasienByID)
		admin.PATCH("/pasien/:id", h.Pasien.UpdatePasien)
		admin.DELETE("/pasien/:id", h.Pasien.DeletePasien)

		admin.GET("/pembayaran/menunggu-konfirmasi", h.Pembayaran.GetMenungguKonfirmasi)
		admin.PUT("/pembayaran/konfirmasi/:id", h.Pembayaran.KonfirmasiPembayaran)
		admin.GET("/pembayaran/tagihan-aktif", h.Pembayaran.GetTagihanAktif)
	}
}
