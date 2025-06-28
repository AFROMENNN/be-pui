package models

import "time"

type RekapNilai struct {
	ID              int       `json:"id"`
	SiswaID         int       `json:"siswa_id"`
	MataPelajaranID int       `json:"mata_pelajaran_id"`
	KelasID         int       `json:"kelas_id"`
	NilaiQuiz       float64   `json:"nilai_quiz"`
	NilaiTugas      float64   `json:"nilai_tugas"`
	NilaiAkhir      float64   `json:"nilai_akhir"`
	Created         time.Time `json:"created"`
	Updated         time.Time `json:"updated"`
}
