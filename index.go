package main

import (
	"be-golang-poliklinik/internal/auth"
	"be-golang-poliklinik/internal/config"
	"be-golang-poliklinik/internal/dokter"
	jadwalperiksa "be-golang-poliklinik/internal/jadwal-periksa"
	"be-golang-poliklinik/internal/obat"
	"be-golang-poliklinik/internal/pasien"
	"be-golang-poliklinik/internal/pembayaran"
	"be-golang-poliklinik/internal/periksa"
	"be-golang-poliklinik/internal/poli"
	"be-golang-poliklinik/internal/riwayat"
	"be-golang-poliklinik/internal/router"
	"flag"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Muat file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment bawaan OS")
	}

	// Atur Mode GIN (Release/Debug)
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	runMigrate := flag.Bool("migrate", false, "Run database migration")
	runSeed := flag.Bool("seed", false, "Run database seeding")
	runFresh := flag.Bool("fresh", false, "Drop tables and run database migration fresh")
	flag.Parse()

	config.ConnectDB()

	if *runFresh {
		config.MigrateFresh()
	} else if *runMigrate {
		config.MigrateDatabase()
	}
	if *runSeed {
		config.SeedAdmin()
	}

	authRepo := auth.NewRepository(config.DB)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	dokterRepo := dokter.NewRepository(config.DB)
	dokterService := dokter.NewService(dokterRepo)
	dokterHandler := dokter.NewHandler(dokterService)

	poliRepo := poli.NewRepository(config.DB)
	poliService := poli.NewService(poliRepo)
	poliHandler := poli.NewHandler(poliService)

	obatRepo := obat.NewRepository(config.DB)
	obatService := obat.NewService(obatRepo)
	obatHandler := obat.NewHandler(obatService)

	pasienRepo := pasien.NewRepository(config.DB)
	pasienService := pasien.NewService(pasienRepo)
	pasienHandler := pasien.NewHandler(pasienService)

	jadwalPeriksaRepo := jadwalperiksa.NewRepository(config.DB)
	jadwalPeriksaService := jadwalperiksa.NewService(jadwalPeriksaRepo)
	jadwalPeriksaHandler := jadwalperiksa.NewHandler(jadwalPeriksaService)

	periksaRepo := periksa.NewRepository(config.DB)
	periksaService := periksa.NewService(periksaRepo)
	periksaHandler := periksa.NewHandler(periksaService)

	riwayatRepo := riwayat.NewRepository(config.DB)
	riwayatService := riwayat.NewService(riwayatRepo)
	riwayatHandler := riwayat.NewHandler(riwayatService)

	pembayaranRepo := pembayaran.NewRepository(config.DB)
	pembayaranService := pembayaran.NewService(pembayaranRepo)
	pembayaranHandler := pembayaran.NewHandler(pembayaranService)

	r := gin.Default()
	
	// Konfigurasi CORS agar aplikasi Flutter Web / Frontend lain bisa memanggil API ini
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	r.MaxMultipartMemory = 5 << 20
	r.Static("/uploads", "./uploads")

	h := &router.Handlers{
		Auth:          authHandler,
		Dokter:        dokterHandler,
		Poli:          poliHandler,
		Obat:          obatHandler,
		Pasien:        pasienHandler,
		JadwalPeriksa: jadwalPeriksaHandler,
		Periksa:       periksaHandler,
		Riwayat:       riwayatHandler,
		Pembayaran:    pembayaranHandler,
	}

	router.SetupRoutes(r, h, authService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server Backend Poliklinik berjalan di port %s...\n", port)
	r.Run(":" + port)
}
