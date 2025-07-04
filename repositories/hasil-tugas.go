package repositories

import (
	"be-pui/models"
	"context"

	"github.com/jmoiron/sqlx"
)

type HasilTugasSiswa struct {
	models.HasilTugas
	NamaSiswa string `db:"nama_siswa"`
}

type HasilTugasKelas struct {
	models.HasilTugas
	NamaSiswa  string `db:"nama_siswa"`
	JudulTugas string `db:"judul_tugas"`
}

type HasilTugasRepository interface {
	Create(ctx context.Context, hasilTugas *models.HasilTugas) error
	GetByTugasAndSiswaID(ctx context.Context, tugasID int, siswaID int) (*models.HasilTugas, error)
	GetAllBySiswaID(ctx context.Context, siswaID int) ([]models.HasilTugas, error)
	GetAllByTugasID(ctx context.Context, tugasID int) ([]HasilTugasSiswa, error)
	GetAllByKelasID(ctx context.Context, kelasID int) ([]HasilTugasKelas, error)
	GetAllByGuruAndMapelID(ctx context.Context, guruID, mapelID int) ([]HasilTugasKelas, error) // Method baru

}

type hasilTugasRepository struct {
	db *sqlx.DB
}

func NewHasilTugasRepository(db *sqlx.DB) HasilTugasRepository {
	return &hasilTugasRepository{db: db}
}

// Create menyisipkan data pengumpulan tugas baru oleh siswa ke dalam database.
func (r *hasilTugasRepository) Create(ctx context.Context, hasilTugas *models.HasilTugas) error {
	query := `
        INSERT INTO hasil_tugas (tugas_id, siswa_id, tanggal_pengumpulan, status, file_jawaban_url)
        VALUES (:tugas_id, :siswa_id, :tanggal_pengumpulan, :status, :file_jawaban_url)
    `
	_, err := r.db.NamedExecContext(ctx, query, hasilTugas)
	return err
}

// GetByTugasAndSiswaID memeriksa apakah seorang siswa sudah mengumpulkan tugas tertentu.
func (r *hasilTugasRepository) GetByTugasAndSiswaID(ctx context.Context, tugasID int, siswaID int) (*models.HasilTugas, error) {
	var hasilTugas models.HasilTugas
	query := "SELECT * FROM hasil_tugas WHERE tugas_id = $1 AND siswa_id = $2"
	err := r.db.GetContext(ctx, &hasilTugas, query, tugasID, siswaID)
	if err != nil {
		return nil, err
	}
	return &hasilTugas, nil
}

// GetAllBySiswaID mengambil semua tugas yang telah dikumpulkan oleh seorang siswa.
func (r *hasilTugasRepository) GetAllBySiswaID(ctx context.Context, siswaID int) ([]models.HasilTugas, error) {
	var hasilTugasList []models.HasilTugas
	query := "SELECT * FROM hasil_tugas WHERE siswa_id = $1 ORDER BY tanggal_pengumpulan DESC"

	err := r.db.SelectContext(ctx, &hasilTugasList, query, siswaID)
	if err != nil {
		return nil, err
	}
	return hasilTugasList, nil
}

// GetAllByTugasID mengambil semua hasil tugas untuk satu tugas tertentu, digabung dengan nama siswa.
func (r *hasilTugasRepository) GetAllByTugasID(ctx context.Context, tugasID int) ([]HasilTugasSiswa, error) {
	var results []HasilTugasSiswa
	query := `
        SELECT
            ht.*,
            s.nama AS nama_siswa
        FROM hasil_tugas ht
        JOIN siswa s ON ht.siswa_id = s.id
        WHERE ht.tugas_id = $1
        ORDER BY s.nama ASC
    `
	err := r.db.SelectContext(ctx, &results, query, tugasID)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// GetAllByKelasID mengambil semua hasil tugas untuk satu kelas tertentu, digabung dengan data siswa dan tugas.
func (r *hasilTugasRepository) GetAllByKelasID(ctx context.Context, kelasID int) ([]HasilTugasKelas, error) {
	var results []HasilTugasKelas
	query := `
		SELECT
			ht.*,
			s.nama AS nama_siswa,
			t.judul AS judul_tugas
		FROM hasil_tugas ht
		JOIN siswa s ON ht.siswa_id = s.id
		JOIN tugas t ON ht.tugas_id = t.id
		WHERE s.kelas_id = $1
		ORDER BY t.deadline DESC, s.nama ASC
	`
	err := r.db.SelectContext(ctx, &results, query, kelasID)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *hasilTugasRepository) GetAllByGuruAndMapelID(ctx context.Context, guruID, mapelID int) ([]HasilTugasKelas, error) {
	var results []HasilTugasKelas
	query := `
		SELECT
			ht.*,
			s.nama AS nama_siswa,
			t.judul AS judul_tugas
		FROM hasil_tugas ht
		JOIN siswa s ON ht.siswa_id = s.id
		JOIN tugas t ON ht.tugas_id = t.id
		JOIN kelas k ON s.kelas_id = k.id
		WHERE k.guru_id = $1 AND t.mata_pelajaran_id = $2
		ORDER BY t.deadline DESC, s.nama ASC
	`
	err := r.db.SelectContext(ctx, &results, query, guruID, mapelID)
	if err != nil {
		return nil, err
	}
	return results, nil
}
