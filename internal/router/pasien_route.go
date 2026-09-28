package router

import (
	"be-golang-poliklinik/internal/auth"
	"be-golang-poliklinik/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetUpPasienRoutes(api *gin.RouterGroup, h *Handlers, authService auth.Service) {
	pasien := api.Group("/pasien")
	pasien.Use(middleware.AuthMiddleware(authService))
	pasien.Use(middleware.RoleBlockMiddleware("pasien"))
	{
		pasien.GET("/jadwal-periksa", h.JadwalPeriksa.GetAllJadwalPeriksa)
		pasien.POST("/booking-periksa", h.Periksa.CreateDaftarPoli)
		pasien.GET("/riwayat", h.Riwayat.GetRiwayatPasien)
		pasien.GET("/riwayat/:id", h.Riwayat.GetRiwayatDetailDokter)
		pasien.POST("/pembayaran", h.Pembayaran.CreatePembayaran)
	}
}
