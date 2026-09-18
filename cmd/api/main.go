package main

import (
	"e-commerce/internal/config"
	"e-commerce/internal/database"
	"e-commerce/internal/product"
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

	// Membuat Repository
	repo := product.NewRepository(pool)

	// Membuat Service
	service := product.NewService(repo)

	// Membuat handler
	handler := product.NewHandler(service)

	// Membuat Router Gin
	router := gin.Default()

	// Endpoint Pengecekan Aplikasi
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Endpoint mengambil semua produk
	router.GET("/products", handler.GetAll)
	router.GET("/products/:id", handler.GetByID)
	router.POST("/products", handler.Create)
	router.PUT("/products/:id", handler.Update)
	router.DELETE("/products/:id", handler.Delete)

	// Menjalankan server
	log.Println("Server running on port", cfg.AppPort)

	err = router.Run(":" + cfg.AppPort)
	if err != nil {
		log.Fatal(err)
	}
}
