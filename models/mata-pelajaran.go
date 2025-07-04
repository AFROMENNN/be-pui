package models

import "time"

type MataPelajaran struct {
	ID        int       `db:"id"`
	Nama      string    `db:"nama"`
	Deskripsi string    `db:"deskripsi"`
	Created   time.Time `db:"created"`
	Updated   time.Time `db:"updated"`
}
