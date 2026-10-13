package cart

import (
	"e-commerce/internal/helper"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(c *gin.Context) {
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

	cart, err := h.service.Create(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "internal server error",
		})

		return
	}

	c.JSON(http.StatusOK, cart)
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

func (h *Handler) AddItem(c *gin.Context) {
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

	var req AddCartItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})

		return
	}

	cartItem, err := h.service.AddItem(c.Request.Context(), id, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})

		return
	}

	c.JSON(http.StatusOK, cartItem)
}

func (h *Handler) GetItems(c *gin.Context) {
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

	items, err := h.service.GetItemsByUserID(
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

		log.Println("GetItems error:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) DeleteCartItem(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user id not found",
		})
		return
	}

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
		})

		return
	}

	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid item ID",
		})

		return
	}

	cart, err := h.service.GetByUserID(c.Request.Context(), uid)
	if err != nil {
		if errors.Is(err, ErrCartNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "cart not found"})
			return
		}

		log.Println("DeleteCartItem find cart error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error"})

		return
	}

	err = h.service.DeleteItem(c.Request.Context(), cart.ID, itemID)

	if err != nil {
		if errors.Is(err, ErrCartItemsNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "cart item not found"})
			return
		}

		log.Println("DeleteCartItem error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error"})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "cart item deleted successfully",
	})

}

func (h *Handler) UpdateItemQuantity(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user id not found",
		})
		return
	}

	var req UpdateCartItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if err := validate.Struct(req); err != nil {
		fieldErrors := helper.ValidationErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "validation failed",
			"fields": fieldErrors,
		})
		return
	}

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
		})
		return
	}

	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid item ID",
		})
		return
	}

	cart, err := h.service.GetByUserID(c.Request.Context(), uid)
	if err != nil {
		if errors.Is(err, ErrCartNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "cart not found",
			})
			return
		}

		log.Println("UpdateItemQuantity find cart error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	CartItem, err := h.service.UpdateItemQuantity(
		c.Request.Context(),
		cart.ID,
		itemID,
		req,
	)
	if err != nil {
		if errors.Is(err, ErrCartItemsNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "cart item not found",
			})
			return
		}

		log.Println("UpdateItemQuantity error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, CartItem)
}
