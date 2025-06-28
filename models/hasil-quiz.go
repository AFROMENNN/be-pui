package models

import "time"

type HasilQuiz struct {
	ID                int       `json:"id"`
	QuizID            int       `json:"quiz_id"`
	SiswaID           int       `json:"siswa_id"`
	Nilai             float64   `json:"nilai"`
	TanggalPengerjaan time.Time `json:"tanggal_pengerjaan"`
	Status            string    `json:"status"`
	Created           time.Time `json:"created"`
	Updated           time.Time `json:"updated"`
}
