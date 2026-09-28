package router

import (
	"be-golang-poliklinik/internal/auth"
	"be-golang-poliklinik/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupDokterRoutes(h *Handlers, api *gin.RouterGroup, authService auth.Service) {
	dokter := api.Group("/dokter")
	dokter.Use(middleware.AuthMiddleware(authService))
	dokter.Use(middleware.RoleBlockMiddleware("dokter"))
	{
		dokter.POST("/jadwal-periksa", h.JadwalPeriksa.CreateJadwalPeriksa)
		dokter.GET("/jadwal-periksa", h.JadwalPeriksa.GetJadwalPeriksa)
		dokter.PUT("/jadwal-periksa/:id", h.JadwalPeriksa.UpdateJadwalPeriksa)
		dokter.DELETE("/jadwal-periksa/:id", h.JadwalPeriksa.DeleteJadwalPeriksa)
		dokter.GET("/obat", h.Obat.GetAllObatAvailable)
		dokter.GET("/daftar-antrian", h.Periksa.GetAntrianDokter)
		dokter.POST("/periksa", h.Periksa.CreatePeriksa)
		dokter.GET("/riwayat", h.Riwayat.GetRiwayatDokter)
		dokter.GET("/riwayat/:id", h.Riwayat.GetRiwayatDetailDokter)
	}
}
