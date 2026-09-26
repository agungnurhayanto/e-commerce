package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, user CreateUserRequest) (*User, error) {
	query := `INSERT INTO users(name, email, password) VALUES ($1, $2, $3) RETURNING id, name, email, password`

	var u User

	err := r.db.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Password,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *Repository) FindAll(ctx context.Context, limit int, offset int) ([]User, error) {
	query := `
			SELECT id, name, email FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var user []User

	for rows.Next() {
		var u User

		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
		)

		if err != nil {
			return nil, err
		}

		user = append(user, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*User, error) {
	query := `SELECT id, name, email, password FROM users WHERE id = $1`

	var u User

	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Password,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &u, nil
}

func (r *Repository) Update(ctx context.Context, id string, req UpdateUserRequest) (*User, error) {
	query := `UPDATE users SET name = $1,email = $2, password= $3, updated_at = NOW() WHERE id = $4 RETURNING id, name, email, password`

	var u User

	err := r.db.QueryRow(
		ctx,
		query,
		req.Name,
		req.Email,
		req.Password,
		id,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Password,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &u, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id=$1::uuid`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return err
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, name, email, password
		FROM users
		WHERE email = $1
	`

	var u User

	err := r.db.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Password,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &u, nil
}
