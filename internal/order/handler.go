package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	// Logika checkout akan kita lengkapi setelah ini.
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Order checkout belum diimplementasikan",
	})
}
