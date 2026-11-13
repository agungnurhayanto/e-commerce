package gateway

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type DummyGateway struct{}

func NewDummyGateway() *DummyGateway {
	return &DummyGateway{}
}

func (d *DummyGateway) CreatePayment(
	ctx context.Context,
	req PaymentRequest,
) (*PaymentResponse, error) {

	transactionID := uuid.New().String()

	return &PaymentResponse{
		TransactionID: transactionID,
		PaymentURL: fmt.Sprintf(
			"http://localhost:8080/dummy-payment/%s",
			transactionID,
		),
		Status: "pending",
	}, nil

}
