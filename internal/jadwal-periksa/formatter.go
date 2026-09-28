package jadwalperiksa

import (
	"be-golang-poliklinik/internal/dokter"
	"be-golang-poliklinik/internal/domain"
	"time"
)

type JadwalPeriksaFormatter struct {
	ID         uint                        `json:"id"`
	Hari       string                      `json:"hari"`
	JamMulai   string                      `json:"jam_mulai"`
	JamSelesai string                      `json:"jam_selesai"`
	Dokter     *dokter.DokterListFormatter `json:"dokter,omitempty"`
}

func FormatJadwal(jadwal domain.JadwalPeriksa) JadwalPeriksaFormatter {
	formatter := JadwalPeriksaFormatter{
		ID:         jadwal.ID,
		Hari:       jadwal.Hari,
		JamMulai:   formatJam(jadwal.JamMulai),
		JamSelesai: formatJam(jadwal.JamSelesai),
	}

	if jadwal.Dokter.ID != 0 {
		dokterFormatted := dokter.FormatDokter(jadwal.Dokter)
		formatter.Dokter = &dokterFormatted
	}

	return formatter
}

func FormatJadwalList(jadwals []domain.JadwalPeriksa) []JadwalPeriksaFormatter {
	var jadwalsFormatter []JadwalPeriksaFormatter

	for _, jadwal := range jadwals {
		jadwalsFormatter = append(jadwalsFormatter, FormatJadwal(jadwal))
	}

	if len(jadwalsFormatter) == 0 {
		return []JadwalPeriksaFormatter{}
	}

	return jadwalsFormatter
}

func formatJam(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("15:04")
}
