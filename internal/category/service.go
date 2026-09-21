package category

import "context"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAll(ctx context.Context, limit int, offset int) ([]Category, error) {
	return s.repository.FindAll(ctx, limit, offset)
}

func (s *Service) GetById(ctx context.Context, id string) (*Category, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, req CreateCategoryRequest) (*Category, error) {
	return s.repository.Create(ctx, req)
}

func (s *Service) Update(ctx context.Context, id string, req UpdateCategoryRequest) (*Category, error) {
	return s.repository.Update(ctx, id, req)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}
