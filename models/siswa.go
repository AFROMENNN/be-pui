package models

import "time"

type Siswa struct {
	ID         int       `json:"id"`
	Nama       string    `json:"nama"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	ProfileUrl string    `json:"profile_url"`
	NoHp       string    `json:"no_hp"`
	KelasID    *int      `json:"kelas_id"` // Nullable foreign key
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}
