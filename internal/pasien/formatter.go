package pasien

import "be-golang-poliklinik/internal/domain"

type PasienFormatter struct {
	ID    uint   `json:"id"`
	Nama  string `json:"nama"`
	Email string `json:"email"`
	NoKTP string `json:"no_ktp"`
	NoRM  string `json:"no_rm"`
}

func FormatPasien(pasien domain.User) PasienFormatter {
	return PasienFormatter{
		ID:    pasien.ID,
		Nama:  pasien.Nama,
		Email: pasien.Email,
		NoKTP: pasien.NoKTP,
		NoRM:  pasien.NoRM,
	}
}

func FormatPasienList(pasiens []domain.User) []PasienFormatter {
	var pasiensFormatter []PasienFormatter

	for _, pasien := range pasiens {
		pasiensFormatter = append(pasiensFormatter, FormatPasien(pasien))
	}

	if len(pasiensFormatter) == 0 {
		return []PasienFormatter{}
	}

	return pasiensFormatter
}
