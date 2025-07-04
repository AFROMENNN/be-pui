package repositories

import (
	"be-pui/models"
	"context"

	"github.com/jmoiron/sqlx"
)

type KelasRepository interface {
	Create(ctx context.Context, kelas *models.Kelas) error
}

type kelasRepository struct {
	db *sqlx.DB
}

func NewKelasRepository(db *sqlx.DB) KelasRepository {
	return &kelasRepository{db: db}
}

func (r *kelasRepository) Create(ctx context.Context, kelas *models.Kelas) error {
	query := `
        INSERT INTO kelas (name, tingkat, guru_id, jumlah_siswa)
        VALUES (:name, :tingkat, :guru_id, :jumlah_siswa)
    `
	_, err := r.db.NamedExecContext(ctx, query, kelas)
	return err
}
