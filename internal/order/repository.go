package order

import (
	"context"
	"errors"

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
		query = `UPDATE products SET stock = stock - $1 WHERE id= $2 AND stock >= $1`
		result, err := tx.Exec(ctx, query, item.Quantity, item.ProductID)

		if err != nil {
			return err
		}

		if result.RowsAffected() == 0 {
			return errors.New("stok produk tidak mencukupi")
		}

		query := `INSERT INTO order_items(order_id, product_id, product_name, unit_price, quantity, subtotal)VALUES($1, $2, $3, $4, $5, $6)`

		_, err = tx.Exec(ctx, query, order.ID, item.ProductID, item.ProductName, item.UnitPrice, item.Quantity, item.Subtotal)
		if err != nil {
			return err
		}
	}

	query = `DELETE FROM cart_items WHERE cart_id = (SELECT id FROM carts WHERE user_id = $1)`
	_, err = tx.Exec(ctx, query, order.UserID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)

}
