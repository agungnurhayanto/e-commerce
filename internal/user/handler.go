package user

import (
	"e-commerce/internal/helper"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var user CreateUserRequest

	// Mengambil data JSON dari request
	if err := c.ShouldBindJSON(&user); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})

		return
	}

	// Memanggil Service dan menampung hasilnya
	createdUser, err := h.service.Create(
		c.Request.Context(), user,
	)

	if err != nil {

		log.Println("Create user error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})

		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

func (h *Handler) GetAll(c *gin.Context) {
	page := 1
	limit := 10

	pageParam := c.Query("page")
	limitParam := c.Query("limit")

	if pageParam != "" {
		parsedPage, err := strconv.Atoi(pageParam)

		if err != nil || parsedPage < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "page must be a positive integer",
			})
			return
		}

		page = parsedPage
	}

	if limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)

		if err != nil || parsedLimit < 1 || parsedLimit > 100 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "limit must be between 1 and 100",
			})
			return
		}

		limit = parsedLimit
	}

	offset := (page - 1) * limit

	users, err := h.service.GetAll(c.Request.Context(), limit, offset)

	if err != nil {
		log.Println("GetAll error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID",
		})

		return
	}

	user, err := h.service.GetById(
		c.Request.Context(),
		id.String(),
	)

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})

			return
		}

		log.Println("GetByID error:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})

		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) Update(c *gin.Context) {

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product ID",
		})

		return
	}

	var req UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if err := validate.Struct(req); err != nil {
		errors := helper.ValidationErrors(err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "validation failed",
			"fields": errors,
		})

		return
	}

	user, err := h.service.Update(c.Request.Context(), id.String(), req)
	if err != nil {

		if errors.Is(err, ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "product not found",
			})

			return
		}

		log.Println("Update error:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product ID",
		})

		return
	}

	err = h.service.Delete(c.Request.Context(), id.String())

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "product not found",
			})
			return
		}

		log.Println("Delete error:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user delete succesfully",
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})

		return
	}

	token, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
