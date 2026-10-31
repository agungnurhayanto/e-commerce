package dto

import "github.com/google/uuid"

type CreatePaymentRequest struct {
	OrderID       uuid.UUID `json:"order_id" validate:"required"`
	PaymentMethod string    `json:"payment_method" validate:"required"`
}
