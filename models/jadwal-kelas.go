package models

import "time"

type JadwalKelas struct {
	ID              int       `json:"id"`
	KelasID         int       `json:"kelas_id"`
	MataPelajaranID int       `json:"mata_pelajaran_id"`
	Hari            string    `json:"hari"`
	JamMulai        time.Time `json:"jam_mulai"`
	JamSelesai      time.Time `json:"jam_selesai"`
	Created         time.Time `json:"created"`
	Updated         time.Time `json:"updated"`
}
