package cart

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID) (*Cart, error) {
	return s.repository.Create(ctx, userID)

}

func (s *Service) GetByUserID(ctx context.Context, id uuid.UUID) (*Cart, error) {
	return s.repository.FindByUserID(ctx, id)
}

func (s *Service) AddItem(ctx context.Context, userID uuid.UUID, req AddCartItemRequest) (*CartItem, error) {
	cart, err := s.repository.FindByUserID(ctx, userID)

	if err != nil {
		return nil, err
	}

	return s.repository.CreateItem(ctx, cart.ID, req)
}

func (s *Service) GetItemsByUserID(ctx context.Context, userID uuid.UUID) ([]CartItem, error) {
	cart, err := s.repository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.repository.FindItemsByCartID(ctx, cart.ID)

}

func (s *Service) DeleteItem(ctx context.Context, cartID, itemID uuid.UUID) error {
	return s.repository.DeleteItem(ctx, cartID, itemID)
}

func (s *Service) UpdateItemQuantity(ctx context.Context, cartID, itemID uuid.UUID, req UpdateCartItemRequest) (*CartItem, error) {
	return s.repository.UpdateItem(ctx, cartID, itemID, req)
}
