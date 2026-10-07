package main

import (
	"log"

	"e-commerce/internal/auth"
	"e-commerce/internal/cart"
	"e-commerce/internal/category"
	"e-commerce/internal/config"
	"e-commerce/internal/database"
	"e-commerce/internal/product"
	"e-commerce/internal/user"

	"github.com/gin-gonic/gin"
)

func main() {

	// CONFIG & DATABASE
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := database.NewPostgresPool(cfg)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer pool.Close()

	// USER
	userRepo := user.NewRepository(pool)
	userService := user.NewService(userRepo, cfg.JWTSecret, cfg.JWTExpireHours)
	userHandler := user.NewHandler(userService)

	// CATEGORY
	categoryRepo := category.NewRepository(pool)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService)

	// PRODUCT
	productRepo := product.NewRepository(pool)
	productService := product.NewService(productRepo, categoryRepo)
	productHandler := product.NewHandler(productService)

	// CART
	cartRepo := cart.NewRepository(pool)
	cartService := cart.NewService(cartRepo)
	cartHandler := cart.NewHandler(cartService)

	// ROUTER
	router := gin.Default()

	// HEALTH
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// USER ROUTES
	router.POST("/login", userHandler.Login)
	router.GET("/me", auth.AuthMiddleware(cfg.JWTSecret), userHandler.GetMe)
	router.POST("/users", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), userHandler.Create)
	router.GET("/users", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), userHandler.GetAll)
	router.GET("/users/:id", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), userHandler.GetByID)
	router.PUT("/users/:id", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), userHandler.Update)
	router.DELETE("/users/:id", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), userHandler.Delete)

	// CATEGORY ROUTES
	router.GET("/categories", categoryHandler.GetAll)
	router.GET("/categories/:id", categoryHandler.GetByID)
	router.POST("/categories", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), categoryHandler.Create)
	router.PUT("/categories/:id", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), categoryHandler.Update)
	router.DELETE("/categories/:id", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), categoryHandler.Delete)

	// PRODUCT ROUTES
	router.GET("/products", productHandler.GetAll)
	router.GET("/products/:id", productHandler.GetByID)
	router.POST("/products", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), productHandler.Create)
	router.PUT("/products/:id", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), productHandler.Update)
	router.DELETE("/products/:id", auth.AuthMiddleware(cfg.JWTSecret), auth.RequireRole("admin"), productHandler.Delete)

	// CART ROUTES
	router.GET("/cart", auth.AuthMiddleware(cfg.JWTSecret), cartHandler.GetByUserID)
	router.POST("/cart", auth.AuthMiddleware(cfg.JWTSecret), cartHandler.Create)

	// SERVER
	log.Println("Server running on port", cfg.AppPort)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
