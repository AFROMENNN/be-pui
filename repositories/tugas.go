package repositories

import (
	"be-pui/models"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type TugasRepository interface {
	Create(ctx context.Context, tugas *models.Tugas) error
	Update(ctx context.Context, tugas *models.Tugas) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*models.Tugas, error)
	GetAll(ctx context.Context) ([]models.Tugas, error)
	GetAllByKelasID(ctx context.Context, kelasID int) ([]models.Tugas, error)
	GetAllByMapelID(ctx context.Context, mapelID int) ([]models.Tugas, error)
	GetAllByKelasAndMapelID(ctx context.Context, kelasID int, mapelID int) ([]models.Tugas, error) // Method baru

}

type tugasRepository struct {
	db *sqlx.DB
}

func NewTugasRepository(db *sqlx.DB) TugasRepository {
	return &tugasRepository{db: db}
}

func (r *tugasRepository) Create(ctx context.Context, tugas *models.Tugas) error {
	query := `
        INSERT INTO tugas (judul, deskripsi, mata_pelajaran_id, kelas_id, deadline)
        VALUES (:judul, :deskripsi, :mata_pelajaran_id, :kelas_id, :deadline)
    `
	_, err := r.db.NamedExecContext(ctx, query, tugas)
	return err
}

func (r *tugasRepository) Update(ctx context.Context, tugas *models.Tugas) error {
	tugas.Updated = time.Now()
	query := `
        UPDATE tugas SET
            judul = :judul,
            deskripsi = :deskripsi,
            mata_pelajaran_id = :mata_pelajaran_id,
            kelas_id = :kelas_id,
            deadline = :deadline,
            updated = :updated
        WHERE id = :id
    `
	_, err := r.db.NamedExecContext(ctx, query, tugas)
	return err
}

func (r *tugasRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM tugas WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *tugasRepository) GetByID(ctx context.Context, id int) (*models.Tugas, error) {
	var tugas models.Tugas
	query := "SELECT * FROM tugas WHERE id = $1"
	err := r.db.GetContext(ctx, &tugas, query, id)
	if err != nil {
		return nil, err
	}
	return &tugas, nil
}

func (r *tugasRepository) GetAll(ctx context.Context) ([]models.Tugas, error) {
	var tugases []models.Tugas
	query := "SELECT * FROM tugas ORDER BY deadline DESC"
	err := r.db.SelectContext(ctx, &tugases, query)
	if err != nil {
		return nil, err
	}
	return tugases, nil
}

// GetAllByKelasID mengambil semua tugas untuk satu kelas tertentu.
func (r *tugasRepository) GetAllByKelasID(ctx context.Context, kelasID int) ([]models.Tugas, error) {
	var tugases []models.Tugas
	query := "SELECT * FROM tugas WHERE kelas_id = $1 ORDER BY deadline DESC"
	err := r.db.SelectContext(ctx, &tugases, query, kelasID)
	if err != nil {
		return nil, err
	}
	return tugases, nil
}

// GetAllByMapelID mengambil semua tugas untuk satu mata pelajaran tertentu.
func (r *tugasRepository) GetAllByMapelID(ctx context.Context, mapelID int) ([]models.Tugas, error) {
	var tugases []models.Tugas
	query := "SELECT * FROM tugas WHERE mata_pelajaran_id = $1 ORDER BY deadline DESC"
	err := r.db.SelectContext(ctx, &tugases, query, mapelID)
	if err != nil {
		return nil, err
	}
	return tugases, nil
}

func (r *tugasRepository) GetAllByKelasAndMapelID(ctx context.Context, kelasID int, mapelID int) ([]models.Tugas, error) {
	var tugases []models.Tugas
	query := "SELECT * FROM tugas WHERE kelas_id = $1 AND mata_pelajaran_id = $2 ORDER BY deadline DESC"
	err := r.db.SelectContext(ctx, &tugases, query, kelasID, mapelID)
	if err != nil {
		return nil, err
	}
	return tugases, nil
}
