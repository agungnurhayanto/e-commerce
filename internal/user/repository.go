package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, user User) (*User, error) {
	query := `INSERT INTO users(name, email, password) VALUES ($1, $2, $3) RETURNING id, name, email, password`

	var u User

	err := r.db.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Password,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
