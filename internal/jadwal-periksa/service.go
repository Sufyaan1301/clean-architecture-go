package jadwalperiksa

import (
	"be-golang-poliklinik/internal/domain"
	"errors"
	"time"
)

type Service interface {
	CreateJadwalPeriksa(idDokter int, input JadwalPeriksaInput) (*domain.JadwalPeriksa, error)
	GetJadwalPeriksaByDokter(idDokter int) ([]domain.JadwalPeriksa, error)
	GetJadwalPeriksaByID(id int) (*domain.JadwalPeriksa, error)
	UpdateJadwalPeriksa(id int, idDokter int, input UpdateJadwalPeriksaInput) (*domain.JadwalPeriksa, error)
	DeleteJadwalPeriksa(id int, idDokter int) error
	GetAllJadwalPeriksa() ([]domain.JadwalPeriksa, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

type JadwalPeriksaInput struct {
	Hari       string `json:"hari" binding:"required"`
	JamMulai   string `json:"jam_mulai" binding:"required"`
	JamSelesai string `json:"jam_selesai" binding:"required"`
}

type UpdateJadwalPeriksaInput struct {
	Hari       string `json:"hari"`
	JamMulai   string `json:"jam_mulai"`
	JamSelesai string `json:"jam_selesai"`
}

func parseJam(jamStr string) (time.Time, error) {
	t, err := time.Parse("15:04", jamStr)
	if err != nil {
		return t, err
	}
	validTime := time.Date(2000, 1, 1, t.Hour(), t.Minute(), 0, 0, time.Local)
	return validTime, nil
}

func (s *service) GetAllJadwalPeriksa() ([]domain.JadwalPeriksa, error) {
	return s.repo.FindAll()
}
func (s *service) CreateJadwalPeriksa(idDokter int, input JadwalPeriksaInput) (*domain.JadwalPeriksa, error) {
	jamMulai, err := parseJam(input.JamMulai)
	if err != nil {
		return nil, errors.New("format jam mulai salah, gunakan format HH:MM (contoh: 08:00)")
	}

	jamSelesai, err := parseJam(input.JamSelesai)
	if err != nil {
		return nil, errors.New("format jam selesai salah, gunakan format HH:MM (contoh: 10:00)")
	}

	if toMinutes(jamSelesai) <= toMinutes(jamMulai) {
		return nil, errors.New("jam selesai harus lebih besar dari jam mulai")
	}

	jadwalHariIni, _ := s.repo.FindByDokterAndHari(idDokter, input.Hari)
	for _, jLama := range jadwalHariIni {
		if isOverlap(jamMulai, jamSelesai, jLama.JamMulai, jLama.JamSelesai) {
			return nil, errors.New("Waktu bertabrakan dengan jadwal Anda yang lain di hari tersebut")
		}
	}

	jadwal := domain.JadwalPeriksa{
		IDDokter:   uint(idDokter),
		Hari:       input.Hari,
		JamMulai:   jamMulai,
		JamSelesai: jamSelesai,
	}

	err = s.repo.Save(&jadwal)
	return &jadwal, err
}

func (s *service) GetJadwalPeriksaByDokter(idDokter int) ([]domain.JadwalPeriksa, error) {
	return s.repo.FindByDokterID(idDokter)
}

func (s *service) GetJadwalPeriksaByID(id int) (*domain.JadwalPeriksa, error) {
	jadwal, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("data jadwal tidak ditemukan")
	}
	return jadwal, nil
}

func (s *service) UpdateJadwalPeriksa(id int, idDokter int, input UpdateJadwalPeriksaInput) (*domain.JadwalPeriksa, error) {
	jadwal, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("data jadwal tidak ditemukan")
	}

	if jadwal.IDDokter != uint(idDokter) {
		return nil, errors.New("akses ditolak: anda tidak diizinkan mengubah jadwal milik dokter lain")
	}

	if input.Hari != "" {
		jadwal.Hari = input.Hari
	}

	if input.JamMulai != "" {
		jamMulai, err := parseJam(input.JamMulai)
		if err != nil {
			return nil, errors.New("format jam mulai salah")
		}
		jadwal.JamMulai = jamMulai
	}

	if input.JamSelesai != "" {
		jamSelesai, err := parseJam(input.JamSelesai)
		if err != nil {
			return nil, errors.New("format jam selesai salah")
		}
		jadwal.JamSelesai = jamSelesai
	}

	if toMinutes(jadwal.JamSelesai) <= toMinutes(jadwal.JamMulai) {
		return nil, errors.New("jam selesai harus lebih besar dari jam mulai")
	}

	jadwalHariIni, _ := s.repo.FindByDokterAndHari(idDokter, jadwal.Hari)
	for _, jLama := range jadwalHariIni {
		if jLama.ID == jadwal.ID {
			continue
		}

		if isOverlap(jadwal.JamMulai, jadwal.JamSelesai, jLama.JamMulai, jLama.JamSelesai) {
			return nil, errors.New("Waktu bertabrakan dengan jadwal Anda yang lain")
		}
	}

	err = s.repo.Update(jadwal)
	return jadwal, err
}

func (s *service) DeleteJadwalPeriksa(id int, idDokter int) error {
	jadwal, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("data jadwal tidak ditemukan")
	}

	if jadwal.IDDokter != uint(idDokter) {
		return errors.New("akses ditolak: anda tidak diizinkan menghapus jadwal milik dokter lain")
	}

	return s.repo.Delete(jadwal)
}

func isOverlap(startBaru, endBaru, startLama, endLama time.Time) bool {
	sb := toMinutes(startBaru)
	eb := toMinutes(endBaru)
	sl := toMinutes(startLama)
	el := toMinutes(endLama)
	return sb < el && eb > sl
}

func toMinutes(t time.Time) int {
	return (t.Hour() * 60) + t.Minute()
}
