package handler

import (
	"be-pui/models"
	"be-pui/repositories"
	"be-pui/utils"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type GuruCreateRequest struct {
	Nama     string `json:"nama" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	NoHp     string `json:"no_hp" binding:"required"`
}

type GuruUpdateRequest struct {
	Nama  string `json:"nama"`
	Email string `json:"email" binding:"omitempty,email"`
	NoHp  string `json:"no_hp"`
}

type GuruResponse struct {
	ID         int       `json:"id"`
	Nama       string    `json:"nama"`
	Email      string    `json:"email"`
	ProfileUrl *string   `json:"profile_url,omitempty"`
	NoHp       string    `json:"no_hp"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

type guruHandler struct {
	guruRepo       repositories.GuruRepository
	tugasRepo      repositories.TugasRepository
	hasilTugasRepo repositories.HasilTugasRepository
	jwtUtil        *utils.JWTUtil
}

type HasilTugasSiswaResponse struct {
	ID                 int       `json:"id"`
	TugasID            int       `json:"tugas_id"`
	SiswaID            int       `json:"siswa_id"`
	NamaSiswa          string    `json:"nama_siswa"`
	JudulTugas         string    `json:"judul_tugas"`
	Nilai              *float64  `json:"nilai,omitempty"`
	TanggalPengumpulan time.Time `json:"tanggal_pengumpulan"`
	Status             string    `json:"status"`
	Feedback           *string   `json:"feedback,omitempty"`
	FileJawabanUrl     *string   `json:"file_jawaban_url,omitempty"`
}

func NewGuruHandler(
	guruRepo repositories.GuruRepository,
	tugasRepo repositories.TugasRepository,
	hasilTugasRepo repositories.HasilTugasRepository,
	jwtUtil *utils.JWTUtil,
) *guruHandler {
	return &guruHandler{
		guruRepo:       guruRepo,
		tugasRepo:      tugasRepo,
		hasilTugasRepo: hasilTugasRepo,
		jwtUtil:        jwtUtil,
	}
}

func (h *guruHandler) CheckTugasSiswa(c *gin.Context) {
	claims, ok := utils.GetCurrentUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Konteks user tidak ditemukan."})
		return
	}
	guruID := claims.UserID

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

	hasilTugas, err := h.hasilTugasRepo.GetAllByGuruAndMapelID(c.Request.Context(), guruID, mapelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil hasil tugas siswa."})
		return
	}

	var response []HasilTugasSiswaResponse
	for _, hasil := range hasilTugas {
		response = append(response, HasilTugasSiswaResponse{
			ID:                 hasil.ID,
			TugasID:            hasil.TugasID,
			SiswaID:            hasil.SiswaID,
			NamaSiswa:          hasil.NamaSiswa,
			JudulTugas:         hasil.JudulTugas,
			Nilai:              hasil.Nilai,
			TanggalPengumpulan: hasil.TanggalPengumpulan,
			Status:             hasil.Status,
			Feedback:           hasil.Feedback,
			FileJawabanUrl:     hasil.FileJawabanUrl,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil hasil tugas siswa.",
		"data":    response,
	})
}

func (h *guruHandler) LoginGuru(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Email atau password tidak valid."})
		return
	}

	guru, err := h.guruRepo.GetByEmail(c.Request.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Kombinasi email dan password salah."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Terjadi kesalahan pada server."})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(guru.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Kombinasi email dan password salah."})
		return
	}

	token, err := h.jwtUtil.GenerateJWTToken(guru)
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

func (h *guruHandler) CreateGuru(c *gin.Context) {
	var req GuruCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Input tidak valid."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal memproses password."})
		return
	}

	guruModel := models.Guru{
		Nama:     req.Nama,
		Email:    req.Email,
		Password: string(hashedPassword),
		NoHp:     req.NoHp,
	}

	if err := h.guruRepo.Create(c.Request.Context(), &guruModel); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Email sudah terdaftar."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Terjadi kesalahan pada server kami."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Guru berhasil dibuat."})
}

func (h *guruHandler) GetProfileGuru(c *gin.Context) {
	claims, ok := utils.GetCurrentUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Konteks user tidak ditemukan. Otentikasi diperlukan.",
		})
		return
	}

	guru, err := h.guruRepo.GetByID(c.Request.Context(), claims.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Profil guru tidak ditemukan."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil profil guru."})
		return
	}

	guruResponse := GuruResponse{
		ID:         guru.ID,
		Nama:       guru.Nama,
		Email:      guru.Email,
		ProfileUrl: guru.ProfileUrl,
		NoHp:       guru.NoHp,
		Created:    guru.Created,
		Updated:    guru.Updated,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil profil guru.",
		"data":    guruResponse,
	})
}

func (h *guruHandler) GetAllGurus(c *gin.Context) {
	gurus, err := h.guruRepo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data guru."})
		return
	}

	var guruResponses []GuruResponse
	for _, guru := range gurus {
		guruResponses = append(guruResponses, GuruResponse{
			ID:         guru.ID,
			Nama:       guru.Nama,
			Email:      guru.Email,
			ProfileUrl: guru.ProfileUrl,
			NoHp:       guru.NoHp,
			Created:    guru.Created,
			Updated:    guru.Updated,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil data seluruh guru.",
		"data":    guruResponses,
	})
}

func (h *guruHandler) GetGuruByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID guru tidak valid."})
		return
	}

	guru, err := h.guruRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Guru tidak ditemukan."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data guru."})
		return
	}

	guruResponse := GuruResponse{
		ID:         guru.ID,
		Nama:       guru.Nama,
		Email:      guru.Email,
		ProfileUrl: guru.ProfileUrl,
		NoHp:       guru.NoHp,
		Created:    guru.Created,
		Updated:    guru.Updated,
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Guru ditemukan.", "data": guruResponse})
}

func (h *guruHandler) UpdateGuru(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID guru tidak valid."})
		return
	}

	var req GuruUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Data yang dikirim tidak valid."})
		return
	}

	existingGuru, err := h.guruRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Guru yang akan diupdate tidak ditemukan."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data guru."})
		return
	}

	if req.Nama != "" {
		existingGuru.Nama = req.Nama
	}
	if req.Email != "" {
		existingGuru.Email = req.Email
	}
	if req.NoHp != "" {
		existingGuru.NoHp = req.NoHp
	}

	if err := h.guruRepo.Update(c.Request.Context(), existingGuru); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal memperbarui data guru."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Data guru berhasil diperbarui."})
}

func (h *guruHandler) DeleteGuru(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID guru tidak valid."})
		return
	}

	_, err = h.guruRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Guru yang akan dihapus tidak ditemukan."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal memeriksa data guru."})
		return
	}

	if err := h.guruRepo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menghapus data guru."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Data guru berhasil dihapus."})
}
