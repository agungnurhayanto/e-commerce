package main

import (
	"e-commerce/internal/category"
	"e-commerce/internal/config"
	"e-commerce/internal/database"
	"e-commerce/internal/product"
	"e-commerce/internal/user"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	// Membaca konfigurasi dari .ent
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

	userService := user.NewService(userRepo)

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

	// Endpoint  User

	router.POST("/users", userHandler.Create)

	// Endpoint  produk
	router.GET("/products", handler.GetAll)
	router.GET("/products/:id", handler.GetByID)
	router.POST("/products", handler.Create)
	router.PUT("/products/:id", handler.Update)
	router.DELETE("/products/:id", handler.Delete)

	// Endpoint category
	router.GET("/categories", categoryHandler.GetAll)
	router.GET("/categories/:id", categoryHandler.GetByID)
	router.POST("/categories", categoryHandler.Create)
	router.PUT("/categories/:id", categoryHandler.Update)
	router.DELETE("/categories/:id", categoryHandler.Delete)

	// Menjalankan server
	log.Println("Server running on port", cfg.AppPort)

	err = router.Run(":" + cfg.AppPort)
	if err != nil {
		log.Fatal(err)
	}
}
