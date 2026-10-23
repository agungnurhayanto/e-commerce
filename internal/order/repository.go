package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrOrderNotFound = errors.New("order not found")

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

func (r *Repository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]Order, error) {
	query := `SELECT id, user_id, total, status, created_at, updated_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []Order

	for rows.Next() {
		var o Order

		err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.Total,
			&o.Status,
			&o.CreatedAt,
			&o.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		itemQuery := `SELECT id, order_id, product_id, product_name, unit_price, quantity, subtotal FROM order_items WHERE order_id=$1`

		itemRows, err := r.db.Query(ctx, itemQuery, o.ID)
		if err != nil {
			return nil, err
		}

		o.Items = make([]OrderItem, 0)

		for itemRows.Next() {
			var item OrderItem

			err := itemRows.Scan(
				&item.ID,
				&item.OrderID,
				&item.ProductID,
				&item.ProductName,
				&item.UnitPrice,
				&item.Quantity,
				&item.Subtotal,
			)

			if err != nil {
				itemRows.Close()
				return nil, err
			}

			o.Items = append(o.Items, item)
		}

		if err := itemRows.Err(); err != nil {
			itemRows.Close()
			return nil, err
		}

		itemRows.Close()

		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
