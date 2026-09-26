package user

import (
	"context"
	"e-commerce/internal/auth"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository     *Repository
	jwtSecret      string
	jwtExpireHours int
}

func NewService(repository *Repository, jwtSecret string,
	jwtExpireHours int) *Service {
	return &Service{
		repository:     repository,
		jwtSecret:      jwtSecret,
		jwtExpireHours: jwtExpireHours,
	}
}

func (s *Service) Create(ctx context.Context, user CreateUserRequest) (*User, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	return s.repository.Create(ctx, user)
}

func (s *Service) GetAll(ctx context.Context, limit int, offset int) ([]User, error) {
	return s.repository.FindAll(ctx, limit, offset)
}

func (s *Service) GetById(ctx context.Context, id string) (*User, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id string, req UpdateUserRequest) (*User, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	req.Password = string(hashedPassword)

	return s.repository.Update(ctx, id, req)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (string, error) {
	user, err := s.repository.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		return "", err
	}

	token, err := auth.GenerateToken(
		user.ID,
		s.jwtSecret,
		s.jwtExpireHours,
	)
	if err != nil {
		return "", err
	}

	return token, nil

}
