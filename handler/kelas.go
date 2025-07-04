package handler

import (
	"be-pui/models"
	"be-pui/repositories"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

type KelasCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Tingkat     int    `json:"tingkat" binding:"required,oneof=4 5 6"`
	JumlahSiswa int    `json:"jumlah_siswa" binding:"required,gte=0"`
	GuruID      int    `json:"guru_id" binding:"required"`
}

type kelasHandler struct {
	kelasRepo repositories.KelasRepository
}

func NewKelasHandler(kelasRepo repositories.KelasRepository) *kelasHandler {
	return &kelasHandler{kelasRepo: kelasRepo}
}

func (h *kelasHandler) CreateKelas(c *gin.Context) {
	var req KelasCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorDetails := make(map[string]string)
			for _, fe := range ve {
				switch fe.Tag() {
				case "required":
					errorDetails[fe.Field()] = "Field ini wajib diisi."
				case "oneof":
					errorDetails[fe.Field()] = "Tingkat harus 4, 5, atau 6."
				case "gte":
					errorDetails[fe.Field()] = "Jumlah siswa tidak boleh negatif."
				default:
					errorDetails[fe.Field()] = "Input tidak valid."
				}
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

	kelasModel := models.Kelas{
		Name:        req.Name,
		Tingkat:     req.Tingkat,
		JumlahSiswa: req.JumlahSiswa,
		GuruID:      req.GuruID,
	}

	if err := h.kelasRepo.Create(c.Request.Context(), &kelasModel); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Guru dengan ID yang diberikan tidak ditemukan.",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server kami.",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Kelas berhasil dibuat.",
		"data":    nil,
	})
}
