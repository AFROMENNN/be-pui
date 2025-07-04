package handler

import (
	"be-pui/config"
	"be-pui/models"
	"be-pui/repositories"
	"be-pui/utils"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type SiswaCreateRequest struct {
	Nama     string `json:"nama" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	NoHp     string `json:"no_hp" binding:"required"`
	KelasID  *int   `json:"kelas_id"`
}

type SiswaResponse struct {
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

type SubmitTugasRequest struct {
	TugasID        int    `json:"tugas_id" binding:"required"`
	FileJawabanUrl string `json:"file_jawaban_url" binding:"required,url"`
}

type TugasWithCompletionStatusResponse struct {
	ID              int                `json:"id"`
	Judul           string             `json:"judul"`
	Deskripsi       string             `json:"deskripsi"`
	MataPelajaranID int                `json:"mata_pelajaran_id"`
	Deadline        time.Time          `json:"deadline"`
	IsCompleted     bool               `json:"is_completed"`
	HasilTugas      *models.HasilTugas `json:"hasil_tugas,omitempty"`
}

type siswaHandler struct {
	siswaRepo      repositories.SiswaRepository
	tugasRepo      repositories.TugasRepository
	hasilTugasRepo repositories.HasilTugasRepository
	jwtUtil        *utils.JWTUtil
	cfg            *config.Config
}

func NewSiswaHandler(
	siswaRepo repositories.SiswaRepository,
	tugasRepo repositories.TugasRepository,
	hasilTugasRepo repositories.HasilTugasRepository,
	jwtUtil *utils.JWTUtil,
	cfg *config.Config,
) *siswaHandler {
	return &siswaHandler{
		siswaRepo:      siswaRepo,
		tugasRepo:      tugasRepo,
		hasilTugasRepo: hasilTugasRepo,
		jwtUtil:        jwtUtil,
		cfg:            cfg,
	}
}

// CheckTugasCompletion memeriksa semua tugas yang ada untuk kelas siswa dan menandai mana yang sudah dikerjakan.
func (h *siswaHandler) CheckTugasCompletion(c *gin.Context) {
	claims, ok := utils.GetCurrentUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Konteks user tidak ditemukan."})
		return
	}

	siswa, err := h.siswaRepo.GetByID(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data profil siswa."})
		return
	}

	if siswa.KelasID == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Siswa belum terdaftar di kelas manapun.",
			"data":    []interface{}{},
		})
		return
	}

	tugasKelas, err := h.tugasRepo.GetAllByKelasID(c.Request.Context(), *siswa.KelasID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil daftar tugas kelas."})
		return
	}

	hasilTugasSiswa, err := h.hasilTugasRepo.GetAllBySiswaID(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil riwayat pengumpulan tugas."})
		return
	}

	hasilMap := make(map[int]models.HasilTugas)
	for _, hasil := range hasilTugasSiswa {
		hasilMap[hasil.TugasID] = hasil
	}

	var response []TugasWithCompletionStatusResponse
	for _, tugas := range tugasKelas {
		tugasItem := TugasWithCompletionStatusResponse{
			ID:              tugas.ID,
			Judul:           tugas.Judul,
			Deskripsi:       tugas.Deskripsi,
			MataPelajaranID: tugas.MataPelajaranID,
			Deadline:        tugas.Deadline,
			IsCompleted:     false,
			HasilTugas:      nil,
		}

		if hasil, found := hasilMap[tugas.ID]; found {
			tugasItem.IsCompleted = true
			tugasItem.HasilTugas = &hasil
		}
		response = append(response, tugasItem)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil status pengerjaan tugas.",
		"data":    response,
	})
}

func (h *siswaHandler) SubmitTugas(c *gin.Context) {
	tugasIDStr := c.PostForm("tugas_id")
	if tugasIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "tugas_id wajib diisi."})
		return
	}
	tugasID, err := strconv.Atoi(tugasIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "tugas_id tidak valid."})
		return
	}

	file, err := c.FormFile("file_jawaban")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "File jawaban wajib di-upload."})
		return
	}

	claims, ok := utils.GetCurrentUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Konteks user tidak ditemukan."})
		return
	}

	_, err = h.hasilTugasRepo.GetByTugasAndSiswaID(c.Request.Context(), tugasID, claims.UserID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Anda sudah pernah mengumpulkan tugas ini."})
		return
	}
	if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal memeriksa status pengumpulan."})
		return
	}

	tugas, err := h.tugasRepo.GetByID(c.Request.Context(), tugasID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Tugas tidak ditemukan."})
		return
	}

	siswa, err := h.siswaRepo.GetByID(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil profil siswa."})
		return
	}

	if siswa.KelasID == nil || *siswa.KelasID != tugas.KelasID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Anda tidak terdaftar di kelas untuk tugas ini."})
		return
	}

	ext := filepath.Ext(file.Filename)
	uniqueFilename := fmt.Sprintf("tugas-%d-siswa-%d-%d%s", tugasID, claims.UserID, time.Now().Unix(), ext)
	dst := filepath.Join("./uploads/jawaban_tugas/", uniqueFilename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan file."})
		return
	}

	status := "selesai"
	if time.Now().After(tugas.Deadline) {
		status = "terlambat"
	}

	baseURL := h.cfg.Server.BaseURL
	fileURL := fmt.Sprintf("%s/static/jawaban_tugas/%s", baseURL, uniqueFilename)

	hasilTugasModel := models.HasilTugas{
		TugasID:            tugasID,
		SiswaID:            claims.UserID,
		TanggalPengumpulan: time.Now(),
		Status:             status,
		FileJawabanUrl:     &fileURL,
	}

	if err := h.hasilTugasRepo.Create(c.Request.Context(), &hasilTugasModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan pengumpulan tugas."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Tugas berhasil dikumpulkan."})
}

func (h *siswaHandler) GetMyTugas(c *gin.Context) {
	mapelIDStr := c.Query("mapel_id")
	if mapelIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Query parameter 'mapel_id' wajib diisi."})
		return
	}
	mapelID, err := strconv.Atoi(mapelIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Query parameter 'mapel_id' tidak valid."})
		return
	}

	claims, ok := utils.GetCurrentUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Konteks user tidak ditemukan."})
		return
	}

	profile, err := h.siswaRepo.GetProfileByID(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data profil siswa."})
		return
	}

	if profile.KelasID == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Siswa belum terdaftar di kelas manapun.",
			"data":    []models.Tugas{},
		})
		return
	}

	tugases, err := h.tugasRepo.GetAllByKelasAndMapelID(c.Request.Context(), *profile.KelasID, mapelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data tugas."})
		return
	}

	var tugasResponses []TugasResponse
	for _, tugas := range tugases {
		tugasResponses = append(tugasResponses, TugasResponse{
			ID:              tugas.ID,
			Judul:           tugas.Judul,
			Deskripsi:       tugas.Deskripsi,
			MataPelajaranID: tugas.MataPelajaranID,
			KelasID:         tugas.KelasID,
			Deadline:        tugas.Deadline,
			Created:         tugas.Created,
			Updated:         tugas.Updated,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar tugas.",
		"data":    tugasResponses,
	})
}

func (h *siswaHandler) LoginSiswa(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Email atau password tidak valid."})
		return
	}

	siswa, err := h.siswaRepo.GetByEmail(c.Request.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Kombinasi email dan password salah."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Terjadi kesalahan pada server."})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(siswa.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Kombinasi email dan password salah."})
		return
	}

	token, err := h.jwtUtil.GenerateJWTToken(siswa)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat token."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
		"data": gin.H{
			"token": token,
		},
	})
}

func (h *siswaHandler) GetProfileSiswa(c *gin.Context) {
	claims, ok := utils.GetCurrentUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Konteks user tidak ditemukan."})
		return
	}

	profile, err := h.siswaRepo.GetProfileByID(c.Request.Context(), claims.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Profil siswa tidak ditemukan."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil profil siswa."})
		return
	}

	siswaResponse := SiswaResponse{
		ID:             profile.ID,
		Nama:           profile.Nama,
		Email:          profile.Email,
		ProfileUrl:     profile.ProfileUrl,
		NoHp:           profile.NoHp,
		KelasID:        profile.KelasID,
		NamaKelas:      profile.NamaKelas,
		NamaWaliKelas:  profile.NamaWaliKelas,
		EmailWaliKelas: profile.EmailWaliKelas,
		NoHpWaliKelas:  profile.NoHpWaliKelas,
		Created:        profile.Created,
		Updated:        profile.Updated,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil profil siswa.",
		"data":    siswaResponse,
	})
}

func (h *siswaHandler) CreateSiswa(c *gin.Context) {
	var req SiswaCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorDetails := make(map[string]string)
			for _, fe := range ve {
				errorDetails[fe.Field()] = "Input tidak valid"
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Terdapat kesalahan pada data yang Anda masukkan.",
				"errors":  errorDetails,
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Request body tidak valid."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal memproses password."})
		return
	}

	siswaModel := models.Siswa{
		Nama:     req.Nama,
		Email:    req.Email,
		Password: string(hashedPassword),
		NoHp:     req.NoHp,
		KelasID:  req.KelasID,
	}

	if err := h.siswaRepo.Create(c.Request.Context(), &siswaModel); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				c.JSON(http.StatusConflict, gin.H{
					"success": false,
					"message": "Email yang Anda masukkan sudah terdaftar.",
				})
				return
			case "23503":
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Kelas dengan ID yang diberikan tidak ditemukan.",
				})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server kami.",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Siswa berhasil dibuat.",
		"data":    nil,
	})
}
