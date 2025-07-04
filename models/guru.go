package models

import "time"

type Guru struct {
	ID         int       `db:"id"`
	Nama       string    `db:"nama"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	ProfileUrl *string   `db:"profile_url"`
	NoHp       string    `db:"no_hp"`
	Created    time.Time `db:"created"`
	Updated    time.Time `db:"updated"`
}

func (a *Guru) GetID() int {
	return a.ID
}

func (a *Guru) GetEmail() string {
	return a.Email
}

func (g *Guru) GetRole() string {
	return "guru"
}
