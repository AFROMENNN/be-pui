package repositories

import (
	"be-pui/models"
	"context"

	"github.com/jmoiron/sqlx"
)

type MapelRepository interface {
	Create(ctx context.Context, mapel *models.MataPelajaran) error
	GetByID(ctx context.Context, id int) (*models.MataPelajaran, error)
	GetAll(ctx context.Context) ([]models.MataPelajaran, error)
}

type mapelRepository struct {
	db *sqlx.DB
}

func NewMapelRepository(db *sqlx.DB) MapelRepository {
	return &mapelRepository{db: db}
}

func (r *mapelRepository) Create(ctx context.Context, mapel *models.MataPelajaran) error {
	query := `
        INSERT INTO mata_pelajaran (nama, deskripsi)
        VALUES (:nama, :deskripsi)
    `
	_, err := r.db.NamedExecContext(ctx, query, mapel)
	return err
}

func (r *mapelRepository) GetByID(ctx context.Context, id int) (*models.MataPelajaran, error) {
	var mapel models.MataPelajaran
	query := "SELECT * FROM mata_pelajaran WHERE id = $1"
	err := r.db.GetContext(ctx, &mapel, query, id)
	if err != nil {
		return nil, err
	}
	return &mapel, nil
}

func (r *mapelRepository) GetAll(ctx context.Context) ([]models.MataPelajaran, error) {
	var mapels []models.MataPelajaran
	query := "SELECT * FROM mata_pelajaran ORDER BY nama ASC"
	err := r.db.SelectContext(ctx, &mapels, query)
	if err != nil {
		return nil, err
	}
	return mapels, nil
}
