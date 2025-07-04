package models

import "time"

type Siswa struct {
	ID         int       `db:"id"`
	Nama       string    `db:"nama"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	ProfileUrl *string   `db:"profile_url"`
	NoHp       string    `db:"no_hp"`
	KelasID    *int      `db:"kelas_id"`
	Created    time.Time `db:"created"`
	Updated    time.Time `db:"updated"`
}

func (a *Siswa) GetID() int {
	return a.ID
}

func (a *Siswa) GetEmail() string {
	return a.Email
}
func (s *Siswa) GetRole() string {
	return "siswa"
}
