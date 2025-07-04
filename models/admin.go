package models

import "time"

type Admin struct {
	ID         int       `db:"id"`
	Nama       string    `db:"nama"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	ProfileUrl *string   `db:"profile_url"`
	NoHp       string    `db:"no_hp"`
	Role       string    `db:"role"`
	Created    time.Time `db:"created"`
	Updated    time.Time `db:"updated"`
}

func (a *Admin) GetID() int {
	return a.ID
}

func (a *Admin) GetEmail() string {
	return a.Email
}

func (a *Admin) GetRole() string {
	return a.Role
}
