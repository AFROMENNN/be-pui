package models

import "time"

type Kelas struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Tingkat     int       `db:"tingkat"`
	JumlahSiswa int       `db:"jumlah_siswa"`
	GuruID      int       `db:"guru_id"`
	Created     time.Time `db:"created"`
	Updated     time.Time `db:"updated"`
}
