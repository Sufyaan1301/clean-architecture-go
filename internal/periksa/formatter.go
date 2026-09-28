package periksa

import (
	"be-golang-poliklinik/internal/domain"
	"be-golang-poliklinik/internal/pasien"
)

type AntrianFormatter struct {
	ID        uint                   `json:"id"`
	NoAntrian int                    `json:"no_antrian"`
	Keluhan   string                 `json:"keluhan"`
	Pasien    pasien.PasienFormatter `json:"pasien"`
}

func FormatAntrian(daftar domain.DaftarPoli) AntrianFormatter {
	formatter := AntrianFormatter{
		ID:        daftar.ID,
		NoAntrian: daftar.NoAntrian,
		Keluhan:   daftar.Keluhan,
	}

	if daftar.Pasien.ID != 0 {
		formatter.Pasien = pasien.FormatPasien(daftar.Pasien)
	}

	return formatter
}

func FormatAntrianList(daftars []domain.DaftarPoli) []AntrianFormatter {
	var formatters []AntrianFormatter

	for _, d := range daftars {
		formatters = append(formatters, FormatAntrian(d))
	}

	if len(formatters) == 0 {
		return []AntrianFormatter{}
	}

	return formatters
}
