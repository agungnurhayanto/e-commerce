package payment

import (
	"context"
	"e-commerce/internal/order"
	"e-commerce/internal/payment/dto"
	"errors"

	paymentgateway "e-commerce/internal/payment/gateway"

	"github.com/google/uuid"
)

type Service struct {
	repository *Repository
	orderRepo  *order.Repository
	gateway    paymentgateway.Gateway
}

func NewService(repository *Repository, orderRepo *order.Repository, gateway paymentgateway.Gateway) *Service {
	return &Service{
		repository: repository,
		orderRepo:  orderRepo,
		gateway:    gateway,
	}
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, req dto.CreatePaymentRequest) (*Payment, error) {

	// 2. Cari order berdasarkan order_id dan user_id
	order, err := s.orderRepo.FindByIDAndUserID(ctx, req.OrderID, userID)
	if err != nil {
		return nil, err
	}

	// 3. Pastikan order masih pending
	if order.Status != "pending" {
		return nil, errors.New("order tidak dapat dibayar")
	}

	// 4. Minta Dummy Gateway membuat transaksi
	gatewayResponse, err := s.gateway.CreatePayment(
		ctx, paymentgateway.PaymentRequest{
			OrderID:       order.ID.String(),
			Amount:        order.Total,
			PaymentMethod: req.PaymentMethod,
		},
	)

	if err != nil {
		return nil, err
	}

	// 5. Simpan transaction reference
	transactionReference := gatewayResponse.TransactionID

	// Buat object payment
	payment := &Payment{
		OrderID:              order.ID,
		PaymentMethod:        req.PaymentMethod,
		Amount:               order.Total,
		Status:               gatewayResponse.Status,
		TransactionReference: &transactionReference,
	}

	err = s.repository.Create(ctx, payment)
	if err != nil {
		return nil, err
	}

	// 6. Kembalikan data payment
	return payment, nil
}

func (s *Service) ConfirmDummyPayment(
	ctx context.Context,
	paymentID uuid.UUID,
	userID uuid.UUID,
) error {
	return s.repository.ConfirmDummyPayment(
		ctx, paymentID, userID,
	)
}
