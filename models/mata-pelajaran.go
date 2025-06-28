package models

import "time"

type MataPelajaran struct {
	ID        int       `json:"id"`
	Nama      string    `json:"nama"`
	Deskripsi string    `json:"deskripsi"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}
