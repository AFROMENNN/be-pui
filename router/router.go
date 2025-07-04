package router

import (
	"be-pui/config"
	"be-pui/handler"
	"be-pui/middleware"
	"be-pui/repositories"
	"be-pui/utils"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func SetupRouter(db *sqlx.DB, cfg *config.Config) *gin.Engine {
	jwtSecret := cfg.SecretKey

	jwtUtil := utils.NewJWTUtil(jwtSecret)
	authMiddleware := middleware.NewAuthMiddleware(jwtUtil)

	// Repositories
	adminRepo := repositories.NewAdminRepository(db)
	guruRepo := repositories.NewGuruRepository(db)
	kelasRepo := repositories.NewKelasRepository(db)
	siswaRepo := repositories.NewSiswaRepository(db)
	mapelRepo := repositories.NewMapelRepository(db)
	tugasRepo := repositories.NewTugasRepository(db)
	hasilTugasRepo := repositories.NewHasilTugasRepository(db)

	// Handlers
	adminHandler := handler.NewAdminHandler(adminRepo, jwtUtil)
	guruHandler := handler.NewGuruHandler(guruRepo, tugasRepo, hasilTugasRepo, jwtUtil)
	kelasHandler := handler.NewKelasHandler(kelasRepo)
	siswaHandler := handler.NewSiswaHandler(siswaRepo, tugasRepo, hasilTugasRepo, jwtUtil, cfg)
	mapelHandler := handler.NewMapelHandler(mapelRepo)
	tugasHandler := handler.NewTugasHandler(tugasRepo)

	router := gin.Default()

	router.MaxMultipartMemory = 8 << 20

	router.Static("/static", "./uploads")

	api := router.Group("/api/v1")
	{
		// --- Rute Admin ---
		adminRoutes := api.Group("/admins")
		{
			adminRoutes.POST("/login", adminHandler.LoginAdmin)

			adminAuthRoutes := adminRoutes.Group("/")
			adminAuthRoutes.Use(authMiddleware.Auth())
			{
				adminGeneralRoutes := adminAuthRoutes.Group("/")
				adminGeneralRoutes.Use(authMiddleware.RequireRole("super admin", "admin biasa"))
				{
					adminGeneralRoutes.GET("/", adminHandler.GetAllAdmins)
					adminGeneralRoutes.GET("/profile", adminHandler.GetProfileAdmin)
					adminGeneralRoutes.GET("/:id", adminHandler.GetAdminByID)
				}

				adminSuperRoutes := adminAuthRoutes.Group("/")
				adminSuperRoutes.Use(authMiddleware.RequireRole("super admin"))
				{
					adminSuperRoutes.POST("/", adminHandler.CreateAdmin)
				}
			}
		}

		// --- Rute Guru ---
		guruRoutes := api.Group("/gurus")
		{
			guruRoutes.POST("/login", guruHandler.LoginGuru)

			guruProfileRoutes := guruRoutes.Group("/")
			guruProfileRoutes.Use(authMiddleware.Auth(), authMiddleware.RequireRole("guru"))
			{
				guruProfileRoutes.GET("/profile", guruHandler.GetProfileGuru)
				guruProfileRoutes.GET("/tugas", guruHandler.CheckTugasSiswa)
			}

			guruManagementRoutes := guruRoutes.Group("/")
			guruManagementRoutes.Use(authMiddleware.Auth(), authMiddleware.RequireRole("super admin", "admin biasa"))
			{
				guruManagementRoutes.POST("/", guruHandler.CreateGuru)
				guruManagementRoutes.GET("/", guruHandler.GetAllGurus)
				guruManagementRoutes.GET("/:id", guruHandler.GetGuruByID)
				guruManagementRoutes.PUT("/:id", guruHandler.UpdateGuru)
				guruManagementRoutes.DELETE("/:id", guruHandler.DeleteGuru)
			}
		}

		// --- Rute Kelas ---
		kelasRoutes := api.Group("/kelas")
		kelasRoutes.Use(authMiddleware.Auth())
		{
			kelasRoutes.POST("/", authMiddleware.RequireRole("super admin", "admin biasa"), kelasHandler.CreateKelas)
		}

		// --- Rute Siswa ---
		siswaRoutes := api.Group("/siswas")
		{
			siswaRoutes.POST("/login", siswaHandler.LoginSiswa)

			siswaProfileRoutes := siswaRoutes.Group("/")
			siswaProfileRoutes.Use(authMiddleware.Auth(), authMiddleware.RequireRole("siswa"))
			{
				siswaProfileRoutes.GET("/profile", siswaHandler.GetProfileSiswa)
				siswaProfileRoutes.GET("/tugas", siswaHandler.GetMyTugas)
				siswaProfileRoutes.POST("/tugas/submit", siswaHandler.SubmitTugas)
				siswaProfileRoutes.GET("/tugas/status", siswaHandler.CheckTugasCompletion)
			}

			siswaManagementRoutes := siswaRoutes.Group("/")
			siswaManagementRoutes.Use(authMiddleware.Auth(), authMiddleware.RequireRole("super admin", "admin biasa"))
			{
				siswaManagementRoutes.POST("/", siswaHandler.CreateSiswa)
				// Rute lain untuk manajemen siswa bisa ditambahkan di sini
			}
		}

		// --- Rute Mata Pelajaran (Mapel) ---
		mapelRoutes := api.Group("/mapels")
		mapelRoutes.Use(authMiddleware.Auth())
		{
			mapelRoutes.POST("/", authMiddleware.RequireRole("super admin", "admin biasa"), mapelHandler.CreateMapel)

			mapelRoutes.GET("/", authMiddleware.RequireRole("super admin", "admin biasa", "guru", "siswa"), mapelHandler.GetAllMapel)
			mapelRoutes.GET("/:id", authMiddleware.RequireRole("super admin", "admin biasa", "guru", "siswa"), mapelHandler.GetByIDMapel)
		}

		// --- Rute Tugas ---
		tugasRoutes := api.Group("/tugas")
		tugasRoutes.Use(authMiddleware.Auth())
		{
			tugasRoutes.POST("/", authMiddleware.RequireRole("guru"), tugasHandler.CreateTugas)
			tugasRoutes.GET("/kelas/:kelas_id", authMiddleware.RequireRole("guru", "siswa", "super admin", "admin biasa"), tugasHandler.GetAllTugasByKelasID)
			tugasRoutes.GET("/mapel/:mapel_id", authMiddleware.RequireRole("guru", "siswa", "super admin", "admin biasa"), tugasHandler.GetAllTugasByMapelID)
		}
	}

	return router
}
