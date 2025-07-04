package models

import "time"

type HasilTugas struct {
	ID                 int       `db:"id"`
	TugasID            int       `db:"tugas_id"`
	SiswaID            int       `db:"siswa_id"`
	Nilai              *float64  `db:"nilai"`
	TanggalPengumpulan time.Time `db:"tanggal_pengumpulan"`
	Status             string    `db:"status"`
	Feedback           *string   `db:"feedback"`
	FileJawabanUrl     *string   `db:"file_jawaban_url"`
	Created            time.Time `db:"created"`
	Updated            time.Time `db:"updated"`
}
