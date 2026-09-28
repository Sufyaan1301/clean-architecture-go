package dokter

import "be-golang-poliklinik/internal/domain"

type DokterListFormatter struct {
	ID       uint   `json:"id"`
	Nama     string `json:"nama"`
	Foto     string `json:"foto"`
	NamaPoli string `json:"nama_poli"`
}

// Format untuk 1 data Dokter (Dipakai oleh modul lain, misal Jadwal)
func FormatDokter(d domain.User) DokterListFormatter {
	namaPoli := ""
	if d.Poli != nil {
		namaPoli = d.Poli.NamaPoli
	}

	return DokterListFormatter{
		ID:       d.ID,
		Nama:     d.Nama,
		Foto:     d.Foto,
		NamaPoli: namaPoli,
	}
}

// Format untuk banyak data Dokter
func FormatDokterList(dokters []domain.User) []DokterListFormatter {
	var formatter []DokterListFormatter

	for _, d := range dokters {
		formatter = append(formatter, FormatDokter(d))
	}
	
	if len(formatter) == 0 {
		return []DokterListFormatter{}
	}
	return formatter
}
