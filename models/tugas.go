package models

import "time"

type Tugas struct {
	ID              int       `db:"id"`
	Judul           string    `db:"judul"`
	Deskripsi       string    `db:"deskripsi"`
	MataPelajaranID int       `db:"mata_pelajaran_id"`
	KelasID         int       `db:"kelas_id"`
	Deadline        time.Time `db:"deadline"`
	Created         time.Time `db:"created"`
	Updated         time.Time `db:"updated"`
}
