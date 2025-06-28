package models

import "time"

type Kelas struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Tingkat     int       `json:"tingkat"`
	JumlahSiswa int       `json:"jumlah_siswa"`
	GuruID      *int      `json:"guru_id"` // Nullable foreign key
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}
