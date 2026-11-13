package payment

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (r *Repository) ConfirmDummyPayment(
	ctx context.Context,
	paymentID uuid.UUID,
	userID uuid.UUID,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	// 1. Ambil payment dan user
	var orderID uuid.UUID
	var paymentStatus string
	var orderStatus string

	query := `SELECT p.order_id, p.status, o.status FROM payments p
	JOIN orders o ON o.id = p.order_id WHERE p.id = $1 AND o.user_id = $2
	FOR UPDATE of p, o`

	err = tx.QueryRow(ctx, query, paymentID, userID).Scan(
		&orderID,
		&paymentStatus,
		&orderStatus,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("payment tidak ditemukan")
		}
		return err
	}

	// 2. Validasi status payment
	if paymentStatus != "pending" {
		return errors.New("payment bukan pending")
	}

	// 3. validasi status order
	if orderStatus != "pending" {
		return errors.New("order bukan pending")
	}

	// 4. Update Payment
	query = `UPDATE payments SET status = 'paid', paid_at = NOW(), updated_at = NOW() WHERE id = $1`
	_, err = tx.Exec(ctx, query, paymentID)
	if err != nil {
		return err
	}

	// 5. Update order
	query = `UPDATE orders SET status = 'paid', updated_at = NOW() WHERE id = $1`
	_, err = tx.Exec(ctx, query, orderID)
	if err != nil {
		return err
	}

	// 6. Commit
	return tx.Commit(ctx)

}
