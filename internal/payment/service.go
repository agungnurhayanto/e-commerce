package payment

import (
	"context"
	"e-commerce/internal/order"
	"e-commerce/internal/payment/dto"
	"errors"

	"github.com/google/uuid"
)

type Service struct {
	repository *Repository
	orderRepo  *order.Repository
}

func NewService(repository *Repository, orderRepo *order.Repository) *Service {
	return &Service{
		repository: repository,
		orderRepo:  orderRepo,
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

	// 4. Buat object Payment
	payment := &Payment{
		OrderID:       order.ID,
		PaymentMethod: req.PaymentMethod,
		Amount:        order.Total,
		Status:        "pending",
	}

	// 5. Simpan ke database
	err = s.repository.Create(ctx, payment)
	if err != nil {
		return nil, err
	}

	// 6. Kembalikan data payment
	return payment, nil
}
