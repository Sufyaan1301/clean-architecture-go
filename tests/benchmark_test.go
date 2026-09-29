package tests

import (
	"testing"
	"golang.org/x/crypto/bcrypt"
)

// BenchmarkBcryptHashing mengukur berapa nanodetik (ns) yang dibutuhkan
// untuk melakukan enkripsi password. Ini berguna untuk memastikan Cost bcrypt
// tidak membuat server menjadi bottleneck (lambat) saat banyak user mendaftar.
func BenchmarkBcryptHashing(b *testing.B) {
	password := []byte("password_rahasia_123")
	b.ResetTimer() // Reset timer agar persiapan data (variabel) tidak dihitung

	for i := 0; i < b.N; i++ {
		bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	}
}

// BenchmarkBcryptCompare mengukur kecepatan validasi login.
// Ini penting karena API Login dipanggil berkali-kali setiap hari.
func BenchmarkBcryptCompare(b *testing.B) {
	password := []byte("password_rahasia_123")
	hashedPassword, _ := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bcrypt.CompareHashAndPassword(hashedPassword, password)
	}
}
