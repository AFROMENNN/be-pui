package handler

import (
	"be-pui/models"
	"be-pui/repositories"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

type TugasCreateRequest struct {
	Judul           string    `json:"judul" binding:"required"`
	Deskripsi       string    `json:"deskripsi"`
	MataPelajaranID int       `json:"mata_pelajaran_id" binding:"required"`
	KelasID         int       `json:"kelas_id" binding:"required"`
	Deadline        time.Time `json:"deadline" binding:"required"`
}

type TugasResponse struct {
	ID              int       `json:"id"`
	Judul           string    `json:"judul"`
	Deskripsi       string    `json:"deskripsi"`
	MataPelajaranID int       `json:"mata_pelajaran_id"`
	KelasID         int       `json:"kelas_id"`
	Deadline        time.Time `json:"deadline"`
	Created         time.Time `json:"created"`
	Updated         time.Time `json:"updated"`
}

type tugasHandler struct {
	tugasRepo repositories.TugasRepository
}

func NewTugasHandler(tugasRepo repositories.TugasRepository) *tugasHandler {
	return &tugasHandler{tugasRepo: tugasRepo}
}

func (h *tugasHandler) CreateTugas(c *gin.Context) {
	var req TugasCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorDetails := make(map[string]string)
			for _, fe := range ve {
				errorDetails[fe.Field()] = "Input tidak valid"
			}
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Terdapat kesalahan pada data yang Anda masukkan.", "errors": errorDetails})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Request body tidak valid."})
		return
	}

	tugasModel := models.Tugas{
		Judul:           req.Judul,
		Deskripsi:       req.Deskripsi,
		MataPelajaranID: req.MataPelajaranID,
		KelasID:         req.KelasID,
		Deadline:        req.Deadline,
	}

	if err := h.tugasRepo.Create(c.Request.Context(), &tugasModel); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Kelas atau Mata Pelajaran dengan ID yang diberikan tidak ditemukan."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Terjadi kesalahan pada server kami."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Tugas berhasil dibuat."})
}

func (h *tugasHandler) GetAllTugasByKelasID(c *gin.Context) {
	kelasID, err := strconv.Atoi(c.Param("kelas_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID kelas tidak valid."})
		return
	}

	tugases, err := h.tugasRepo.GetAllByKelasID(c.Request.Context(), kelasID)
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
		"message": "Berhasil mengambil data tugas untuk kelas ini.",
		"data":    tugasResponses,
	})
}

func (h *tugasHandler) GetAllTugasByMapelID(c *gin.Context) {
	mapelID, err := strconv.Atoi(c.Param("mapel_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID mata pelajaran tidak valid."})
		return
	}

	tugases, err := h.tugasRepo.GetAllByMapelID(c.Request.Context(), mapelID)
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
		"message": "Berhasil mengambil data tugas untuk mata pelajaran ini.",
		"data":    tugasResponses,
	})
}
