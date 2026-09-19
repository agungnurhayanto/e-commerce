package product

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrProductNotFound = errors.New("product not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) FindAll(ctx context.Context) ([]Product, error) {
	query := `
			SELECT id, name, description, price, stock FROM products ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Stock,
		)

		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*Product, error) {
	query := `SELECT id, name, description, price, stock FROM products WHERE id = $1`

	var p Product

	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}

		return nil, err
	}

	return &p, nil
}

func (r *Repository) Create(ctx context.Context, req CreateProductRequest) (*Product, error) {
	query := `INSERT INTO products(name, description, price, stock) VALUES($1, $2, $3, $4) RETURNING id,name,description,price,stock`

	var p Product

	err := r.db.QueryRow(
		ctx,
		query,
		req.Name,
		req.Description,
		req.Price,
		req.Stock,
	).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *Repository) Update(ctx context.Context, id string, req UpdateProductRequest) (*Product, error) {
	query := `UPDATE products SET name = $1,description = $2, price = $3, stock = $4, updated_at = NOW() WHERE id = $5 RETURNING id, name, description, price, stock`

	var p Product

	err := r.db.QueryRow(
		ctx,
		query,
		req.Name,
		req.Description,
		req.Price,
		req.Stock,
		id,
	).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}

		return nil, err
	}

	return &p, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id=$1::uuid`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrProductNotFound
	}

	return err
}
