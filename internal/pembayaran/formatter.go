package pembayaran

import (
	"be-golang-poliklinik/internal/domain"
	"time"
)

type TagihanFormatter struct {
	IDPeriksa    uint   `json:"id_periksa"`
	TglPeriksa   string `json:"tgl_periksa"`
	IDPasien     uint   `json:"id_pasien"`
	NamaPasien   string `json:"nama_pasien"`
	TotalTagihan int    `json:"total_tagihan"`
}

func FormatTagihan(p domain.Periksa) TagihanFormatter {
	namaPasien := ""
	var idPasien uint = 0

	if p.DaftarPoli != nil {
		idPasien = p.DaftarPoli.IDPasien
		namaPasien = p.DaftarPoli.Pasien.Nama
	}

	return TagihanFormatter{
		IDPeriksa:    p.ID,
		TglPeriksa:   formatTgl(p.TglPeriksa),
		IDPasien:     idPasien,
		NamaPasien:   namaPasien,
		TotalTagihan: p.BiayaPeriksa,
	}
}

func FormatTagihanList(periksas []domain.Periksa) []TagihanFormatter {
	var formatters []TagihanFormatter
	for _, p := range periksas {
		formatters = append(formatters, FormatTagihan(p))
	}
	if len(formatters) == 0 {
		return []TagihanFormatter{}
	}
	return formatters
}

type PembayaranFormatter struct {
	IDPembayaran uint   `json:"id_pembayaran"`
	IDPeriksa    uint   `json:"id_periksa"`
	TglBayar     string `json:"tgl_bayar"`
	JumlahBayar  int    `json:"jumlah_bayar"`
	MetodeBayar  string `json:"metode_bayar"`
	Status       string `json:"status"`
}

func FormatPembayaran(b domain.Pembayaran) PembayaranFormatter {
	return PembayaranFormatter{
		IDPembayaran: b.ID,
		IDPeriksa:    b.IDPeriksa,
		TglBayar:     formatTgl(b.TglBayar),
		JumlahBayar:  b.JumlahBayar,
		MetodeBayar:  b.MetodeBayar,
		Status:       b.Status,
	}
}

func FormatPembayaranList(pembayarans []domain.Pembayaran) []PembayaranFormatter {
	var formatters []PembayaranFormatter
	for _, b := range pembayarans {
		formatters = append(formatters, FormatPembayaran(b))
	}
	if len(formatters) == 0 {
		return []PembayaranFormatter{}
	}
	return formatters
}

func formatTgl(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02 Jan 2006 15:04")
}
