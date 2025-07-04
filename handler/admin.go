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
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type AdminCreateRequest struct {
	Nama     string `json:"nama" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	NoHp     string `json:"no_hp" binding:"required"`
}

type AdminResponse struct {
	ID         int       `json:"id"`
	Nama       string    `json:"nama"`
	Email      string    `json:"email"`
	ProfileUrl *string   `json:"profile_url,omitempty"`
	NoHp       string    `json:"no_hp"`
	Role       string    `json:"role"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type adminHandler struct {
	adminRepo repositories.AdminRepository
	jwtUtil   *utils.JWTUtil
}

func NewAdminHandler(adminRepo repositories.AdminRepository, jwtUtil *utils.JWTUtil) *adminHandler {
	return &adminHandler{adminRepo: adminRepo, jwtUtil: jwtUtil}
}

func (h *adminHandler) LoginAdmin(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Email atau password tidak valid.",
			"data":    nil,
			"errors":  nil,
		})
		return
	}

	admin, err := h.adminRepo.GetByEmail(c.Request.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Kombinasi email dan password salah.",
				"data":    nil,
				"errors":  nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Terjadi kesalahan pada server."})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Kombinasi email dan password salah.",
			"data":    nil,
			"errors":  nil,
		})
		return
	}

	token, err := h.jwtUtil.GenerateJWTToken(admin)
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
		"errors": nil,
	})
}

func (h *adminHandler) CreateAdmin(c *gin.Context) {
	var req AdminCreateRequest

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
				"data":    nil,
				"errors":  errorDetails,
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Request body tidak valid.",
			"data":    nil,
			"errors":  nil,
		})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	adminModel := models.Admin{
		Nama:     req.Nama,
		Email:    req.Email,
		Password: string(hashedPassword),
		NoHp:     req.NoHp,
		Role:     "admin biasa",
	}

	if err := h.adminRepo.Create(c.Request.Context(), &adminModel); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Email yang Anda masukkan sudah terdaftar. Silakan gunakan email lain.",
				"data":    nil,
				"errors":  nil,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server kami. Mohon coba lagi.",
			"data":    nil,
			"errors":  nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Admin berhasil dibuat",
		"data":    nil,
		"errors":  nil,
	})
}

func (h *adminHandler) GetAllAdmins(c *gin.Context) {
	admins, err := h.adminRepo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data admin.",
			"data":    nil,
			"errors":  nil,
		})
		return
	}

	var adminResponses []AdminResponse
	for _, admin := range admins {
		adminResponses = append(adminResponses, AdminResponse{
			ID:         admin.ID,
			Nama:       admin.Nama,
			Email:      admin.Email,
			ProfileUrl: admin.ProfileUrl,
			NoHp:       admin.NoHp,
			Role:       admin.Role,
			Created:    admin.Created,
			Updated:    admin.Updated,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil data seluruh admin.",
		"data":    adminResponses,
		"errors":  nil,
	})
}

func (h *adminHandler) GetAdminByID(c *gin.Context) {
	id := c.Param("id")

	adminID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid.",
			"data":    nil,
			"errors":  nil,
		})
		return
	}

	admin, err := h.adminRepo.GetByID(c.Request.Context(), adminID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Admin tidak ditemukan.",
				"data":    nil,
				"errors":  nil,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server.",
			"data":    nil,
			"errors":  nil,
		})
		return
	}

	adminResponse := AdminResponse{
		ID:         admin.ID,
		Nama:       admin.Nama,
		Email:      admin.Email,
		ProfileUrl: admin.ProfileUrl,
		NoHp:       admin.NoHp,
		Role:       admin.Role,
		Created:    admin.Created,
		Updated:    admin.Updated,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Admin ditemukan.",
		"data":    adminResponse,
		"errors":  nil,
	})
}

func (h *adminHandler) GetProfileAdmin(c *gin.Context) {
	claims, ok := utils.GetCurrentUserClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Konteks user tidak ditemukan. Otentikasi diperlukan.",
			"data":    nil,
			"errors":  nil,
		})
		return
	}

	admin, err := h.adminRepo.GetByID(c.Request.Context(), claims.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Profil admin tidak ditemukan.",
				"data":    nil,
				"errors":  nil,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil profil admin.",
			"data":    nil,
			"errors":  nil,
		})
		return
	}

	adminResponse := AdminResponse{
		ID:         admin.ID,
		Nama:       admin.Nama,
		Email:      admin.Email,
		ProfileUrl: admin.ProfileUrl,
		NoHp:       admin.NoHp,
		Role:       admin.Role,
		Created:    admin.Created,
		Updated:    admin.Updated,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil profil admin.",
		"data":    adminResponse,
		"errors":  nil,
	})
}
