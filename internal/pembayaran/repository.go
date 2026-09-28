package pembayaran

import (
	"be-golang-poliklinik/internal/domain"

	"gorm.io/gorm"
)

type Repository interface {
	GetTagihanAktif() ([]domain.Periksa, error)
	GetMenungguKonfirmasi() ([]domain.Pembayaran, error)
	GetPeriksaByID(idPeriksa int) (*domain.Periksa, error)
	SavePembayaran(pembayaran *domain.Pembayaran) error
	KonfirmasiPembayaran(idPembayaran int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetTagihanAktif() ([]domain.Periksa, error) {
	var tagihans []domain.Periksa
	err := r.db.
		Preload("DaftarPoli").
		Preload("DaftarPoli.Pasien").
		Joins("LEFT JOIN pembayarans ON pembayarans.id_periksa = periksas.id").
		Where("pembayarans.id IS NULL").
		Order("periksas.tgl_periksa ASC").
		Find(&tagihans).Error
	return tagihans, err
}

func (r *repository) GetPeriksaByID(idPeriksa int) (*domain.Periksa, error) {
	var periksa domain.Periksa
	err := r.db.
		Joins("LEFT JOIN pembayarans ON pembayarans.id_periksa = periksas.id").
		Where("pembayarans.id IS NULL").
		Where("periksas.id = ?", idPeriksa).
		First(&periksa).Error
	return &periksa, err
}

func (r *repository) SavePembayaran(pembayaran *domain.Pembayaran) error {
	return r.db.Create(pembayaran).Error
}

func (r *repository) GetMenungguKonfirmasi() ([]domain.Pembayaran, error) {
	var pembayarans []domain.Pembayaran
	err := r.db.
		Preload("Periksa.DaftarPoli.Pasien").
		Where("status = ?", "Menunggu Konfirmasi").
		Order("tgl_bayar ASC").
		Find(&pembayarans).Error
	return pembayarans, err
}

func (r *repository) KonfirmasiPembayaran(idPembayaran int) error {
	return r.db.Model(&domain.Pembayaran{}).
		Where("id = ?", idPembayaran).
		Update("status", "Lunas").Error
}
