package order

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Total     float64
	Status    string
	Items     []OrderItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderItem struct {
	ID          uuid.UUID
	OrderID     uuid.UUID
	ProductID   uuid.UUID
	ProductName string
	UnitPrice   float64
	Quantity    int
	Subtotal    float64
}
