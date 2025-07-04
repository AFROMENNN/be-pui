package repositories

import (
	"be-pui/models"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type GuruRepository interface {
	Create(ctx context.Context, guru *models.Guru) error
	Update(ctx context.Context, guru *models.Guru) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*models.Guru, error)
	GetByEmail(ctx context.Context, email string) (*models.Guru, error)
	GetAll(ctx context.Context) ([]models.Guru, error)
}

type guruRepository struct {
	db *sqlx.DB
}

func NewGuruRepository(db *sqlx.DB) GuruRepository {
	return &guruRepository{db: db}
}

func (r *guruRepository) Create(ctx context.Context, guru *models.Guru) error {
	query := `
        INSERT INTO guru (nama, email, no_hp, password, profile_url)
        VALUES (:nama, :email, :no_hp, :password, :profile_url)
    `
	_, err := r.db.NamedExecContext(ctx, query, guru)
	return err
}

func (r *guruRepository) Update(ctx context.Context, guru *models.Guru) error {
	guru.Updated = time.Now()
	query := `
        UPDATE guru SET
            nama = :nama,
            email = :email,
            no_hp = :no_hp,
            password = :password,
            profile_url = :profile_url,
            updated = :updated
        WHERE id = :id
    `
	_, err := r.db.NamedExecContext(ctx, query, guru)
	return err
}

func (r *guruRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM guru WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *guruRepository) GetByID(ctx context.Context, id int) (*models.Guru, error) {
	var guru models.Guru
	query := "SELECT * FROM guru WHERE id = $1"
	err := r.db.GetContext(ctx, &guru, query, id)
	if err != nil {
		return nil, err
	}
	return &guru, nil
}

func (r *guruRepository) GetByEmail(ctx context.Context, email string) (*models.Guru, error) {
	var guru models.Guru
	query := "SELECT * FROM guru WHERE email = $1"
	err := r.db.GetContext(ctx, &guru, query, email)
	if err != nil {
		return nil, err
	}
	return &guru, nil
}

func (r *guruRepository) GetAll(ctx context.Context) ([]models.Guru, error) {
	var gurus []models.Guru
	query := "SELECT * FROM guru ORDER BY id DESC"
	err := r.db.SelectContext(ctx, &gurus, query)
	if err != nil {
		return nil, err
	}
	return gurus, nil
}
