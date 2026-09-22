package product

import (
	"context"
	"e-commerce/internal/category"
)

type Service struct {
	repository         *Repository
	categoryRepository *category.Repository
}

func NewService(
	repository *Repository,
	categoryRepository *category.Repository,
) *Service {
	return &Service{
		repository:         repository,
		categoryRepository: categoryRepository,
	}
}

func (s *Service) GetAll(ctx context.Context, limit int, offset int) ([]Product, error) {
	return s.repository.FindAll(ctx, limit, offset)
}

func (s *Service) GetById(ctx context.Context, id string) (*Product, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, req CreateProductRequest) (*Product, error) {
	if req.CategoryID != nil {
		_, err := s.categoryRepository.FindActiveByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, err
		}
	}

	return s.repository.Create(ctx, req)

}

func (s *Service) Update(ctx context.Context, id string, req UpdateProductRequest) (*Product, error) {
	if req.CategoryID != nil {
		_, err := s.categoryRepository.FindActiveByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, err
		}
	}

	return s.repository.Update(ctx, id, req)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}
