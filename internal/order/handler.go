package order

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User tidak terauntetikasi"})

		return
	}

	userIDString, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Tipe user ID tidak valid",
		})
		return
	}

	uid, err := uuid.Parse(userIDString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Format user ID tidak valid",
		})
		return
	}

	err = h.service.Create(c.Request.Context(), uid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal membuat order",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "order berhasil dibuat",
	})
}

func (h *Handler) GetByUserID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User tidak terauntetikasi"})

		return
	}

	userIDString, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Tipe user ID tidak valid",
		})
		return
	}

	uid, err := uuid.Parse(userIDString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Format user ID tidak valid",
		})
		return
	}

	orders, err := h.service.GetByUserID(
		c.Request.Context(),
		uid,
	)

	if err != nil {
		log.Println("GetByUserID error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})

		return
	}
	if orders == nil {
		orders = []Order{}
	}

	c.JSON(http.StatusOK, orders)
}

func (h *Handler) UpdateStatus(c *gin.Context) {

	// ambil ID order dari url
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID order tidak valid",
		})

		return
	}
	// Ambil status dari JSON body
	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Request body tidak valid",
		})

		return
	}

	// Panggil service
	err = h.service.UpdateStatus(c.Request.Context(), orderID, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status order berhasil di perbarui",
	})

}
