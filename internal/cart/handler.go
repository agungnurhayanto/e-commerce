package cart

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetByUserID(c *gin.Context) {
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user id not found",
		})

		return
	}

	id, err := uuid.Parse(userID.(string))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
		})
		return
	}

	cart, err := h.service.GetByUserID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		if errors.Is(err, ErrCartNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "cart not found",
			})

			return
		}

		log.Println("GetByUserID error:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})

		return
	}

	c.JSON(http.StatusOK, cart)
}
