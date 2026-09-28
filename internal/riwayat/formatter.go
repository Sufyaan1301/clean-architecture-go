package riwayat

import (
	"be-golang-poliklinik/internal/domain"
	"fmt"
	"time"
)

type RiwayatDokterFormatter struct {
	ID         uint     `json:"id"`
	Tgl        string   `json:"tgl_periksa"`
	NamaPasien string   `json:"nama_pasien"`
	Keluhan    string   `json:"keluhan"`
	Catatan    string   `json:"catatan"`
	Obat       []string `json:"obat"`
}

func FormatRiwayatDokter(p domain.Periksa) RiwayatDokterFormatter {
	var obats []string
	for _, d := range p.DetailPeriksa {
		obats = append(obats, fmt.Sprintf("%s (%dx)", d.Obat.NamaObat, d.Jumlah))
	}

	namaPasien := ""
	keluhan := ""
	if p.DaftarPoli != nil {
		namaPasien = p.DaftarPoli.Pasien.Nama
		keluhan = p.DaftarPoli.Keluhan
	}

	return RiwayatDokterFormatter{
		ID:         p.ID,
		Tgl:        formatTgl(p.TglPeriksa),
		NamaPasien: namaPasien,
		Keluhan:    keluhan,
		Catatan:    p.Catatan,
		Obat:       obats,
	}
}

func FormatRiwayatDokterList(periksas []domain.Periksa) []RiwayatDokterFormatter {
	var formatters []RiwayatDokterFormatter
	for _, p := range periksas {
		formatters = append(formatters, FormatRiwayatDokter(p))
	}
	if len(formatters) == 0 {
		return []RiwayatDokterFormatter{}
	}
	return formatters
}

type ListRiwayatPasienFormatter struct {
	ID          uint     `json:"id"`
	Tgl         string   `json:"tgl_periksa"`
	NamaDokter  string   `json:"nama_dokter"`
	NamaPoli    string   `json:"nama_poli"`
	Keluhan     string   `json:"keluhan"`
	Catatan     string   `json:"catatan"`
	TotalBiaya  int      `json:"total_biaya"`
	StatusBayar string   `json:"status_bayar"`
	Obat        []string `json:"obat"`
}

func FormatRiwayatPasien(p domain.Periksa) ListRiwayatPasienFormatter {
	var obats []string
	for _, d := range p.DetailPeriksa {
		obats = append(obats, fmt.Sprintf("%s (%dx)", d.Obat.NamaObat, d.Jumlah))
	}

	namaDokter := ""
	namaPoli := ""
	keluhan := ""
	if p.DaftarPoli != nil {
		keluhan = p.DaftarPoli.Keluhan
		namaDokter = p.DaftarPoli.Jadwal.Dokter.Nama
		if p.DaftarPoli.Jadwal.Dokter.Poli != nil {
			namaPoli = p.DaftarPoli.Jadwal.Dokter.Poli.NamaPoli
		}
	}

	statusBayar := "Belum Dibayar"
	if p.Pembayaran != nil {
		statusBayar = "Lunas"
	}

	return ListRiwayatPasienFormatter{
		ID:          p.ID,
		Tgl:         formatTgl(p.TglPeriksa),
		NamaDokter:  namaDokter,
		NamaPoli:    namaPoli,
		Keluhan:     keluhan,
		Catatan:     p.Catatan,
		TotalBiaya:  p.BiayaPeriksa,
		StatusBayar: statusBayar,
		Obat:        obats,
	}
}

func FormatRiwayatPasienList(periksas []domain.Periksa) []ListRiwayatPasienFormatter {
	var formatters []ListRiwayatPasienFormatter
	for _, p := range periksas {
		formatters = append(formatters, FormatRiwayatPasien(p))
	}
	if len(formatters) == 0 {
		return []ListRiwayatPasienFormatter{}
	}
	return formatters
}

func formatTgl(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02 Jan 2006 15:04") 
}
