package gateway

import "context"

type PaymentRequest struct {
	OrderID       string
	Amount        float64
	PaymentMethod string
}

type PaymentResponse struct {
	TransactionID string
	PaymentURL    string
	Status        string
}

type Gateway interface {
	CreatePayment(
		ctx context.Context,
		req PaymentRequest,
	) (*PaymentResponse, error)
}
