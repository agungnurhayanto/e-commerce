package dto

import "time"

type PaymentResponse struct {
	ID                   string     `json:"id"`
	OrderID              string     `json:"order_id"`
	PaymentMethod        string     `json:"payment_method"`
	Amount               int64      `json:"amount"`
	Status               string     `json:"status"`
	TransactionReference string     `json:"transaction_reference"`
	PaidAt               *time.Time `json:"paid_at"`
	CreatedAt            time.Time  `json:"created_at"`
}
