package domain

import (
	"time"
)

type Poli struct {
	ID         uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	NamaPoli   string `json:"nama_poli" gorm:"type:varchar(255);not null"`
	Keterangan string `json:"keterangan" gorm:"type:text"`

	// Relasi: Satu Poli punya banyak Dokter (Users)
	Dokter []User `json:"dokter,omitempty" gorm:"foreignKey:IDPoli"`
}

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleDokter Role = "dokter"
	RolePasien Role = "pasien"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Nama     string `json:"nama" gorm:"type:varchar(255);not null"`
	Alamat   string `json:"alamat" gorm:"type:varchar(255)"`
	NoKTP    string `json:"no_ktp" gorm:"type:varchar(255);unique"`
	NoHP     string `json:"no_hp" gorm:"type:varchar(255)"`
	NoRM     string `json:"no_rm" gorm:"type:varchar(255)"` // Nomor Rekam Medis (untuk pasien)
	Role     Role   `json:"role" gorm:"type:enum('admin','dokter','pasien');not null"`
	Email    string `json:"email" gorm:"type:varchar(255);unique;not null"`
	Password string `json:"-" gorm:"type:varchar(255);not null"`     // "-" disembunyikan dari JSON
	Foto     string `json:"foto,omitempty" gorm:"type:varchar(255)"` // Path ke file foto

	// Foreign Key untuk Dokter (nullable, karena pasien/admin tidak punya poli)
	IDPoli *uint `json:"id_poli"`
	Poli   *Poli `json:"poli,omitempty" gorm:"foreignKey:IDPoli"` // Relasi BelongsTo

	// Relasi
	JadwalPeriksa []JadwalPeriksa `json:"jadwal_periksa,omitempty" gorm:"foreignKey:IDDokter"`
	DaftarPoli    []DaftarPoli    `json:"daftar_poli,omitempty" gorm:"foreignKey:IDPasien"`
}

type JadwalPeriksa struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDDokter   uint      `json:"id_dokter" gorm:"not null"`
	Dokter     User      `json:"dokter,omitempty" gorm:"foreignKey:IDDokter"`
	Hari       string    `json:"hari" gorm:"type:enum('Senin','Selasa','Rabu','Kamis','Jumat','Sabtu','Minggu');not null"`
	JamMulai   time.Time `json:"jam_mulai" gorm:"type:time;not null"`
	JamSelesai time.Time `json:"jam_selesai" gorm:"type:time;not null"`

	// Relasi
	DaftarPoli []DaftarPoli `json:"daftar_poli,omitempty" gorm:"foreignKey:IDJadwal"`
}

type DaftarPoli struct {
	ID        uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	IDPasien  uint          `json:"id_pasien" gorm:"not null"`
	Pasien    User          `json:"pasien,omitempty" gorm:"foreignKey:IDPasien"`
	IDJadwal  uint          `json:"id_jadwal" gorm:"not null"`
	Jadwal    JadwalPeriksa `json:"jadwal,omitempty" gorm:"foreignKey:IDJadwal"`
	Keluhan   string        `json:"keluhan" gorm:"type:text;not null"`
	NoAntrian int           `json:"no_antrian" gorm:"not null"`

	// Relasi
	Periksa *Periksa `json:"periksa,omitempty" gorm:"foreignKey:IDDaftarPoli"` // Relasi HasOne
}

type Periksa struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDDaftarPoli uint      `json:"id_daftar_poli" gorm:"not null;unique"` // Unique menjadikannya One-to-One
	DaftarPoli   *DaftarPoli `json:"daftar_poli,omitempty" gorm:"foreignKey:IDDaftarPoli"`
	TglPeriksa   time.Time `json:"tgl_periksa" gorm:"type:datetime;not null"`
	Catatan      string    `json:"catatan" gorm:"type:text;not null"`
	BiayaPeriksa int       `json:"biaya_periksa" gorm:"not null"`

	// Relasi
	DetailPeriksa []DetailPeriksa `json:"detail_periksa,omitempty" gorm:"foreignKey:IDPeriksa"`
	Pembayaran    *Pembayaran     `json:"pembayaran,omitempty" gorm:"foreignKey:IDPeriksa"`
}

type Obat struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	NamaObat string `json:"nama_obat" gorm:"type:varchar(255);not null"`
	Kemasan  string `json:"kemasan" gorm:"type:varchar(255)"`
	Harga    int    `json:"harga" gorm:"not null"`
	Stok     int    `json:"stok" gorm:"not null"`

	// Relasi
	DetailPeriksa []DetailPeriksa `json:"detail_periksa,omitempty" gorm:"foreignKey:IDObat"`
}

type DetailPeriksa struct {
	ID        uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	IDPeriksa uint    `json:"id_periksa" gorm:"not null"`
	Periksa   Periksa `json:"periksa,omitempty" gorm:"foreignKey:IDPeriksa"`
	IDObat    uint    `json:"id_obat" gorm:"not null"`
	Obat      Obat    `json:"obat,omitempty" gorm:"foreignKey:IDObat"`
	Jumlah    int     `json:"jumlah" gorm:"not null;default:1"`
}

type BlacklistToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Token     string    `gorm:"type:varchar(500);unique;not null"`
	ExpiredAt time.Time `gorm:"not null"`
}

type Pembayaran struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	IDPeriksa   uint      `json:"id_periksa" gorm:"not null"`
	Periksa     Periksa   `json:"periksa" gorm:"foreignKey:IDPeriksa"`
	TglBayar    time.Time `json:"tgl_bayar"`
	JumlahBayar int       `json:"jumlah_bayar"`
	MetodeBayar string    `json:"metode_bayar" gorm:"type:enum('Tunai','Transfer Bank','QRIS');default:'Tunai'"`
	Status      string    `json:"status" gorm:"type:enum('Menunggu Konfirmasi','Lunas');default:'Menunggu Konfirmasi'"`
}
