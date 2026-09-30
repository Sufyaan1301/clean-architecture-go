package config

import (
	"be-golang-poliklinik/internal/domain"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Fallback jika .env belum dimuat dengan benar (berguna untuk testing / backward compatibility)
	if dbUser == "" {
		dbUser = "root"
		dbHost = "127.0.0.1"
		dbPort = "3306"
		dbName = "poliklinik_golang_db"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}
}

func SeedAdmin() {
	var count int64
	DB.Model(&domain.User{}).Where("role = ?", domain.RoleAdmin).Count(&count)

	if count == 0 {
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.MinCost)

		admin := domain.User{
			Nama:     "Super Admin Poliklinik",
			Email:    "admin@poliklinik.com",
			Password: string(passwordHash),
			Role:     domain.RoleAdmin,
			NoKTP:    "0000000000000000",
		}

		if err := DB.Create(&admin).Error; err != nil {
			log.Println("Gagal menjalankan seeder admin:", err)
		} else {
			log.Println("Akun Admin berhasil dibuat! (admin@poliklinik.com / admin123)")
		}
	}
}

func MigrateDatabase() {
	err := DB.AutoMigrate(
		&domain.Poli{},
		&domain.User{},
		&domain.JadwalPeriksa{},
		&domain.DaftarPoli{},
		&domain.Periksa{},
		&domain.Obat{},
		&domain.DetailPeriksa{},
		&domain.BlacklistToken{},
		&domain.Pembayaran{}, // <-- Ini yang ketinggalan!
	)
	if err != nil {
		log.Fatal("Gagal migrasi:", err)
	}
	log.Println("Migrasi ERD Poliklinik Berhasil!")
}

func MigrateFresh() {
	log.Println("Membuang (Drop) semua tabel di database...")
	
	err := DB.Migrator().DropTable(
		&domain.Pembayaran{}, // <-- Tambahkan di sini (hapus paling pertama karena dia nyantol ke periksa)
		&domain.DetailPeriksa{},
		&domain.Periksa{},
		&domain.DaftarPoli{},
		&domain.JadwalPeriksa{},
		&domain.Obat{},
		&domain.User{},
		&domain.Poli{},
		&domain.BlacklistToken{},
	)
	
	if err != nil {
		log.Fatal("Gagal menghapus tabel:", err)
	}

	log.Println("Tabel berhasil dihapus! Menjalankan ulang migrasi...")
	MigrateDatabase()
}
