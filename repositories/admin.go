package repositories

import (
	"be-pui/models"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type AdminRepository interface {
	Create(ctx context.Context, admin *models.Admin) error
	Update(ctx context.Context, admin *models.Admin) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*models.Admin, error)
	GetAll(ctx context.Context) ([]models.Admin, error)
	UpdateProfileURL(ctx context.Context, id int, profileURL *string) error
	GetByEmail(ctx context.Context, email string) (*models.Admin, error)
}

type adminRepository struct {
	db *sqlx.DB
}

func NewAdminRepository(db *sqlx.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) Create(ctx context.Context, admin *models.Admin) error {
	query := `
        INSERT INTO admin (nama, email, password, profile_url, no_hp, role)
        VALUES (:nama, :email, :password, :profile_url, :no_hp, :role)
    `
	_, err := r.db.NamedExecContext(ctx, query, admin)
	return err
}

func (r *adminRepository) Update(ctx context.Context, admin *models.Admin) error {
	admin.Updated = time.Now()
	query := `
        UPDATE admin SET -- FIX: Menggunakan tabel "admin"
            nama = :nama,
            email = :email,
            password = :password,
            profile_url = :profile_url,
            no_hp = :no_hp,
            role = :role,
			updated = :updated
        WHERE id = :id
    `
	_, err := r.db.NamedExecContext(ctx, query, admin)
	return err
}

func (r *adminRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM admin WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *adminRepository) GetByID(ctx context.Context, id int) (*models.Admin, error) {
	var admin models.Admin
	query := "SELECT * FROM admin WHERE id = $1"

	err := r.db.GetContext(ctx, &admin, query, id)
	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func (r *adminRepository) GetAll(ctx context.Context) ([]models.Admin, error) {
	var admins []models.Admin
	query := "SELECT * FROM admin ORDER BY id DESC"

	err := r.db.SelectContext(ctx, &admins, query)
	if err != nil {
		return nil, err
	}

	return admins, nil
}

func (r *adminRepository) UpdateProfileURL(ctx context.Context, id int, profileURL *string) error {
	query := "UPDATE admin SET profile_url = $1, updated = $2 WHERE id = $3"
	_, err := r.db.ExecContext(ctx, query, profileURL, time.Now(), id)
	return err
}

func (r *adminRepository) GetByEmail(ctx context.Context, email string) (*models.Admin, error) {
	var admin models.Admin
	query := "SELECT * FROM admin WHERE email = $1"

	err := r.db.GetContext(ctx, &admin, query, email)
	if err != nil {
		return nil, err
	}

	return &admin, nil
}
