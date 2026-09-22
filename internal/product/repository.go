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

func (r *Repository) FindAll(ctx context.Context, limit int, offset int) ([]Product, error) {

	query := `select p.id, p.name, p.description, p.price, p.stock, p.category_id, c.name as category_name FROM products p 
	LEFT JOIN categories c ON p.category_id = c.id ORDER BY p.created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []Product
	var categoryName *string

	for rows.Next() {
		var p Product

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.CategoryID,
			&categoryName,
		)

		if err != nil {
			return nil, err
		}

		if p.CategoryID != nil && categoryName != nil {
			p.Category = &CategorySummary{
				ID:   *p.CategoryID,
				Name: *categoryName,
			}
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*Product, error) {
	query := `
        SELECT
            p.id,
            p.name,
            p.description,
            p.price,
            p.stock,
            p.category_id,
            c.name AS category_name
        FROM products p
        LEFT JOIN categories c ON p.category_id = c.id
        WHERE p.id = $1
    `

	var p Product
	var categoryName *string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
		&p.CategoryID,
		&categoryName,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}

		return nil, err
	}

	if p.CategoryID != nil && categoryName != nil {
		p.Category = &CategorySummary{
			ID:   *p.CategoryID,
			Name: *categoryName,
		}
	}

	return &p, nil
}

func (r *Repository) Create(ctx context.Context, req CreateProductRequest) (*Product, error) {
	query := `INSERT INTO products(name, description, price, stock, category_id) VALUES($1, $2, $3, $4, $5) RETURNING id,name,description,price,stock,category_id`

	var p Product

	err := r.db.QueryRow(
		ctx,
		query,
		req.Name,
		req.Description,
		req.Price,
		req.Stock,
		req.CategoryID,
	).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
		&p.CategoryID,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *Repository) Update(ctx context.Context, id string, req UpdateProductRequest) (*Product, error) {
	query := `UPDATE products SET name = $1,description = $2, price = $3, stock = $4, category_id = $5, updated_at = NOW() WHERE id = $6 RETURNING id, name, description, price, stock, category_id`

	var p Product

	err := r.db.QueryRow(
		ctx,
		query,
		req.Name,
		req.Description,
		req.Price,
		req.Stock,
		req.CategoryID,
		id,
	).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
		&p.CategoryID,
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
