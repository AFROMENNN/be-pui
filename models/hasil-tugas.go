package models

import "time"

type HasilTugas struct {
	ID                 int       `json:"id"`
	TugasID            int       `json:"tugas_id"`
	SiswaID            int       `json:"siswa_id"`
	Nilai              float64   `json:"nilai"`
	TanggalPengumpulan time.Time `json:"tanggal_pengumpulan"`
	Status             string    `json:"status"`
	Feedback           string    `json:"feedback"`
	Created            time.Time `json:"created"`
	Updated            time.Time `json:"updated"`
}
