package main

import (
	"e-commerce/internal/auth"
	"e-commerce/internal/category"
	"e-commerce/internal/config"
	"e-commerce/internal/database"
	"e-commerce/internal/product"
	"e-commerce/internal/user"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	// Membaca konfigurasi dari .env
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Membuat Koneksi PostgreSQL
	pool, err := database.NewPostgresPool(cfg)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	defer pool.Close()

	userRepo := user.NewRepository(pool)

	userService := user.NewService(userRepo, cfg.JWTSecret,
		cfg.JWTExpireHours)

	userHandler := user.NewHandler(userService)

	// Membuat Repository
	repo := product.NewRepository(pool)

	// Membuat Repository Category
	categoryRepo := category.NewRepository(pool)

	// Membuat Service
	service := product.NewService(repo, categoryRepo)

	// Membuat handler
	handler := product.NewHandler(service)

	// Membuat Service Category
	categoryService := category.NewService(categoryRepo)

	// Membuat Handler Category
	categoryHandler := category.NewHandler(categoryService)

	// Membuat Router Gin
	router := gin.Default()

	// Endpoint Pengecekan Aplikasi
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Login Aplikasi
	router.POST("/login", userHandler.Login)

	// Endpoint  User
	router.GET(
		"/me",
		auth.AuthMiddleware(cfg.JWTSecret),
		userHandler.GetMe,
	)
	router.POST("/users", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), userHandler.Create)
	router.GET("/users", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), userHandler.GetAll)
	router.PUT("/users/:id", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), userHandler.Update)
	router.GET("/users/:id", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), userHandler.GetByID)
	router.DELETE("/users/:id", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), userHandler.Delete)

	// Endpoint  produk
	router.GET("/products", handler.GetAll)
	router.GET("/products/:id", handler.GetByID)

	router.POST("/products", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), handler.Create)
	router.PUT("/products/:id", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), handler.Update)
	router.DELETE("/products/:id", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), handler.Delete)

	// Endpoint category
	router.GET("/categories", categoryHandler.GetAll)
	router.GET("/categories/:id", categoryHandler.GetByID)

	router.POST("/categories", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), categoryHandler.Create)
	router.PUT("/categories/:id", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), categoryHandler.Update)
	router.DELETE("/categories/:id", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), categoryHandler.Delete)

	// Menjalankan server
	log.Println("Server running on port", cfg.AppPort)

	err = router.Run(":" + cfg.AppPort)
	if err != nil {
		log.Fatal(err)
	}
}
