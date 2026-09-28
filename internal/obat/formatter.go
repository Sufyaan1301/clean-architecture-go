package obat

import (
	"be-golang-poliklinik/internal/domain"
)

type ObatAvailableFormatter struct {
	ID      uint   `json:"id"`
	Nama    string `json:"nama"`
	Kemasan string `json:"kemasan"`
	Stok    int    `json:"stok"`
}

// Format satuan (1 obat)
func FormatObatAvailable(obat domain.Obat) ObatAvailableFormatter {
	return ObatAvailableFormatter{
		ID:      obat.ID,
		Nama:    obat.NamaObat,
		Kemasan: obat.Kemasan,
		Stok:    obat.Stok,
	}
}

// Format banyak (list obat)
func FormatObatAvailableList(obats []domain.Obat) []ObatAvailableFormatter {
	var obatsFormatter []ObatAvailableFormatter

	for _, obat := range obats {
		obatsFormatter = append(obatsFormatter, FormatObatAvailable(obat))
	}

	if len(obatsFormatter) == 0 {
		return []ObatAvailableFormatter{}
	}

	return obatsFormatter
}
