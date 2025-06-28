package models

import "time"

type Quiz struct {
	ID              int       `json:"id"`
	Judul           string    `json:"judul"`
	Deskripsi       string    `json:"deskripsi"`
	MataPelajaranID int       `json:"mata_pelajaran_id"`
	KelasID         int       `json:"kelas_id"`
	Created         time.Time `json:"created"`
	Updated         time.Time `json:"updated"`
}
