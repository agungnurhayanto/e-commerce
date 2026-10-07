package cart

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrCartNotFound = errors.New("cart not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, userID uuid.UUID) (*Cart, error) {
	query := `INSERT INTO carts(user_id) VALUES($1) RETURNING id,user_id,created_at,updated_at`

	var c Cart

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&c.ID,
		&c.UserID,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *Repository) FindByUserID(ctx context.Context, id uuid.UUID) (*Cart, error) {
	query := `SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id = $1`

	var u Cart

	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.UserID,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCartNotFound
		}

		return nil, err
	}

	return &u, nil
}
