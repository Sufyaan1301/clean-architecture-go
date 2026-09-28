package poli

import (
	"be-golang-poliklinik/internal/dokter"
	"be-golang-poliklinik/internal/domain"
)

type PoliListFormatter struct {
	ID         uint   `json:"id"`
	Nama       string `json:"nama"`
	Keterangan string `json:"keterangan"`
}

type DetailPoliFormatter struct {
	ID         uint                         `json:"id"`
	Nama       string                       `json:"nama"`
	Keterangan string                       `json:"keterangan"`
	Dokters    []dokter.DokterListFormatter `json:"dokters"`
}


func FormatPoliList(polis []domain.Poli) []PoliListFormatter {
	var formatter []PoliListFormatter

	for _, d := range polis {
		poliFormatter := PoliListFormatter{
			ID:         d.ID,
			Nama:       d.NamaPoli,
			Keterangan: d.Keterangan,
		}

		formatter = append(formatter, poliFormatter)
	}
	if len(formatter) == 0 {
		return []PoliListFormatter{}
	}
	return formatter
}

func FormatPoliDetail(p domain.Poli) DetailPoliFormatter {
	return DetailPoliFormatter{
		ID:         p.ID,
		Nama:       p.NamaPoli,
		Keterangan: p.Keterangan,
		Dokters:    dokter.FormatDokterList(p.Dokter), 
	}
}
