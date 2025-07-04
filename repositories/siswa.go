package repositories

import (
	"be-pui/models"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type SiswaProfile struct {
	ID             int       `db:"id"`
	Nama           string    `db:"nama"`
	Email          string    `db:"email"`
	ProfileUrl     *string   `db:"profile_url"`
	NoHp           string    `db:"no_hp"`
	KelasID        *int      `db:"kelas_id"`
	NamaKelas      *string   `db:"nama_kelas"`
	NamaWaliKelas  *string   `db:"nama_wali_kelas"`
	EmailWaliKelas *string   `db:"email_wali_kelas"`
	NoHpWaliKelas  *string   `db:"no_hp_wali_kelas"`
	Created        time.Time `db:"created"`
	Updated        time.Time `db:"updated"`
}

type SiswaRepository interface {
	Create(ctx context.Context, siswa *models.Siswa) error
	GetByID(ctx context.Context, id int) (*models.Siswa, error)
	GetByEmail(ctx context.Context, email string) (*models.Siswa, error)
	GetProfileByID(ctx context.Context, id int) (*SiswaProfile, error)
}

type siswaRepository struct {
	db *sqlx.DB
}

func NewSiswaRepository(db *sqlx.DB) SiswaRepository {
	return &siswaRepository{db: db}
}

func (r *siswaRepository) Create(ctx context.Context, siswa *models.Siswa) error {
	query := `
        INSERT INTO siswa (nama, email, no_hp, password, profile_url, kelas_id)
        VALUES (:nama, :email, :no_hp, :password, :profile_url, :kelas_id)
    `
	_, err := r.db.NamedExecContext(ctx, query, siswa)
	return err
}

func (r *siswaRepository) GetByID(ctx context.Context, id int) (*models.Siswa, error) {
	var siswa models.Siswa
	query := "SELECT * FROM siswa WHERE id = $1"
	err := r.db.GetContext(ctx, &siswa, query, id)
	if err != nil {
		return nil, err
	}
	return &siswa, nil
}

func (r *siswaRepository) GetByEmail(ctx context.Context, email string) (*models.Siswa, error) {
	var siswa models.Siswa
	query := "SELECT * FROM siswa WHERE email = $1"
	err := r.db.GetContext(ctx, &siswa, query, email)
	if err != nil {
		return nil, err
	}
	return &siswa, nil
}

func (r *siswaRepository) GetProfileByID(ctx context.Context, id int) (*SiswaProfile, error) {
	var profile SiswaProfile
	query := `
        SELECT
            s.id,
            s.nama,
            s.email,
            s.profile_url,
            s.no_hp,
            s.kelas_id,
            s.created,
            s.updated,
            k.name AS nama_kelas,
            g.nama AS nama_wali_kelas,
            g.email AS email_wali_kelas,
            g.no_hp AS no_hp_wali_kelas
        FROM siswa s
        LEFT JOIN kelas k ON s.kelas_id = k.id
        LEFT JOIN guru g ON k.guru_id = g.id
        WHERE s.id = $1
    `
	err := r.db.GetContext(ctx, &profile, query, id)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}
