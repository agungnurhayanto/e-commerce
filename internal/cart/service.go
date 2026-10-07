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
