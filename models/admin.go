package models

import "time"

type Admin struct {
	ID         int       `json:"id"`
	Nama       string    `json:"nama"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	ProfileUrl string    `json:"profile_url"`
	NoHp       string    `json:"no_hp"`
	Role       string    `json:"role"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}
