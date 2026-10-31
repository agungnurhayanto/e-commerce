package payment

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

func (r *Repository) Create(ctx context.Context, payment *Payment) error {

	query := `INSERT INTO payments (
		order_id,
		payment_method,
		amount,
		status,
		transaction_reference
	)VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		ctx,
		query,
		payment.OrderID,
		payment.PaymentMethod,
		payment.Amount,
		payment.Status,
		payment.TransactionReference,
	).Scan(
		&payment.ID,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil

}
