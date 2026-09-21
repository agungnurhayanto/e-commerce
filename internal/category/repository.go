package category

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrCategoryNotFound = errors.New("category not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) FindAll(ctx context.Context, limit int, offset int) ([]Category, error) {
	query := `
			SELECT id, name, description FROM categories ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var c Category

		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
		)

		if err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*Category, error) {
	query := `SELECT id, name, description FROM categories WHERE id = $1`

	var c Category

	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID,
		&c.Name,
		&c.Description,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}

		return nil, err
	}

	return &c, nil
}

func (r *Repository) Create(ctx context.Context, req CreateCategoryRequest) (*Category, error) {
	query := `INSERT INTO categories(name, description) VALUES($1, $2) RETURNING id,name,description`

	var c Category

	err := r.db.QueryRow(
		ctx,
		query,
		req.Name,
		req.Description,
	).Scan(
		&c.ID,
		&c.Name,
		&c.Description,
	)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *Repository) Update(ctx context.Context, id string, req UpdateCategoryRequest) (*Category, error) {
	query := `UPDATE categories SET name = $1,description = $2, updated_at = NOW() WHERE id = $3 RETURNING id, name, description`

	var c Category

	err := r.db.QueryRow(
		ctx,
		query,
		req.Name,
		req.Description,
		id,
	).Scan(
		&c.ID,
		&c.Name,
		&c.Description,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}

		return nil, err
	}

	return &c, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM categories WHERE id=$1::uuid`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCategoryNotFound
	}

	return err
}
