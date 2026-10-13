package order

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

func (r *Repository) Create(ctx context.Context, order *Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	query := `INSERT INTO orders (user_id, total, status) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	err = tx.QueryRow(ctx, query, order.UserID, order.Total, order.Status).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		query := `INSERT INTO order_items(order_id, product_id, product_name, unit_price, quantity, subtotal)VALUES($1, $2, $3, $4, $5, $6)`

		_, err = tx.Exec(ctx, query, order.ID, item.ProductID, item.ProductName, item.Quantity, item.Subtotal)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)

}
