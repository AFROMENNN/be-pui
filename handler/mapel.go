package handler

import (
	"be-pui/models"
	"be-pui/repositories"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MapelCreateRequest struct {
	Nama      string `json:"nama" binding:"required"`
	Deskripsi string `json:"deskripsi"`
}

type MapelResponse struct {
	ID        int       `json:"id"`
	Nama      string    `json:"nama"`
	Deskripsi string    `json:"deskripsi,omitempty"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}

type mapelHandler struct {
	mapelRepo repositories.MapelRepository
}

func NewMapelHandler(mapelRepo repositories.MapelRepository) *mapelHandler {
	return &mapelHandler{mapelRepo: mapelRepo}
}

func (h *mapelHandler) CreateMapel(c *gin.Context) {
	var req MapelCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Input tidak valid."})
		return
	}

	mapelModel := models.MataPelajaran{
		Nama:      req.Nama,
		Deskripsi: req.Deskripsi,
	}

	if err := h.mapelRepo.Create(c.Request.Context(), &mapelModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan mata pelajaran."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Mata pelajaran berhasil dibuat."})
}

func (h *mapelHandler) GetAllMapel(c *gin.Context) {
	mapels, err := h.mapelRepo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data mata pelajaran."})
		return
	}

	var mapelResponses []MapelResponse
	for _, mapel := range mapels {
		mapelResponses = append(mapelResponses, MapelResponse{
			ID:        mapel.ID,
			Nama:      mapel.Nama,
			Deskripsi: mapel.Deskripsi,
			Created:   mapel.Created,
			Updated:   mapel.Updated,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil data seluruh mata pelajaran.",
		"data":    mapelResponses,
	})
}

func (h *mapelHandler) GetByIDMapel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID mata pelajaran tidak valid."})
		return
	}

	mapel, err := h.mapelRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Mata pelajaran tidak ditemukan."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data mata pelajaran."})
		return
	}

	mapelResponse := MapelResponse{
		ID:        mapel.ID,
		Nama:      mapel.Nama,
		Deskripsi: mapel.Deskripsi,
		Created:   mapel.Created,
		Updated:   mapel.Updated,
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Mata pelajaran ditemukan.", "data": mapelResponse})
}
