package cart

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrCartNotFound = errors.New("cart not found")

var ErrCartItemsNotFound = errors.New("cart item not found")

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

func (r *Repository) CreateItem(ctx context.Context, cartID uuid.UUID, req AddCartItemRequest) (*CartItem, error) {
	query := `INSERT INTO cart_items(cart_id, product_id, quantity) VALUES($1, $2, $3) RETURNING id,cart_id,product_id, quantity`

	var c CartItem

	err := r.db.QueryRow(
		ctx,
		query,
		cartID,
		req.ProductID,
		req.Quantity,
	).Scan(
		&c.ID,
		&c.CartID,
		&c.ProductID,
		&c.Quantity,
	)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *Repository) FindItemsByCartID(
	ctx context.Context,
	cartID uuid.UUID,
) ([]CartItem, error) {
	query := `SELECT id, cart_id,product_id,quantity, created_at, updated_at FROM cart_items WHERE cart_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, cartID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []CartItem

	for rows.Next() {
		var item CartItem

		err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.ProductID,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil

}

func (r *Repository) DeleteItem(ctx context.Context, cartID uuid.UUID, itemID uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE id = $1 AND cart_id= $2`

	result, err := r.db.Exec(ctx, query, itemID, cartID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCartItemsNotFound
	}

	return err
}

func (r *Repository) UpdateItem(ctx context.Context, cartID uuid.UUID, itemID uuid.UUID, req UpdateCartItemRequest) (*CartItem, error) {
	query := `UPDATE cart_items SET quantity = $1,updated_at = NOW() WHERE id = $2 and cart_id= $3 RETURNING id, quantity`

	var c CartItem

	err := r.db.QueryRow(
		ctx,
		query,
		req.Quantity,
		itemID,
		cartID,
	).Scan(
		&c.ID,
		&c.Quantity,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCartItemsNotFound
		}

		return nil, err
	}

	return &c, nil
}
