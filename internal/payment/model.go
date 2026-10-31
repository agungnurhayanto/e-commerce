package payment

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID                   uuid.UUID  `json:"id"`
	OrderID              uuid.UUID  `json:"order_id"`
	PaymentMethod        string     `json:"payment_method"`
	Amount               float64    `json:"amount"`
	Status               string     `json:"status"`
	TransactionReference *string    `json:"transaction_reference"`
	PaidAt               *time.Time `json:"paid_at"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}
