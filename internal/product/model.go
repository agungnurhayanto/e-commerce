package product

import "github.com/google/uuid"

type Product struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Stock       int        `json:"stock"`
	CategoryID  *uuid.UUID `json:"category_id"`
}

type CreateProductRequest struct {
	Name        string     `json:"name" validate:"required,min=3,max=100"`
	Description string     `json:"description" validate:"max=1000"`
	Price       float64    `json:"price" validate:"gte=0"`
	Stock       int        `json:"stock" validate:"gte=0"`
	CategoryID  *uuid.UUID `json:"category_id"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Description string  `json:"description" validate:"max=1000"`
	Price       float64 `json:"price" validate:"gte=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
}
