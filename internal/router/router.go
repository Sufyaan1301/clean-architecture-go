package router

import (
	"be-golang-poliklinik/internal/auth"
	"be-golang-poliklinik/internal/dokter"
	jadwalperiksa "be-golang-poliklinik/internal/jadwal-periksa"
	"be-golang-poliklinik/internal/obat"
	"be-golang-poliklinik/internal/pasien"
	"be-golang-poliklinik/internal/pembayaran"
	"be-golang-poliklinik/internal/periksa"
	"be-golang-poliklinik/internal/poli"
	"be-golang-poliklinik/internal/riwayat"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth          *auth.AuthHandler
	Dokter        *dokter.DokterHandler
	Poli          *poli.PoliHandler
	Obat          *obat.ObatHandler
	Pasien        *pasien.PasienHandler
	JadwalPeriksa *jadwalperiksa.JadwalPeriksaHandler
	Periksa       *periksa.PeriksaHandler
	Riwayat       *riwayat.RiwayatHandler
	Pembayaran    *pembayaran.PembayaranHandler
}

func SetupRoutes(r *gin.Engine, h *Handlers, authService auth.Service) {
	api := r.Group("/api")

	setupAuthRoutes(api, h, authService)
	setupAdminRoutes(api, h, authService)
	SetupDokterRoutes(h, api, authService)
	SetUpPasienRoutes(api, h, authService)
}
