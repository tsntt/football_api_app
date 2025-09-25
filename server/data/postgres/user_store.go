package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/tsntt/footballapi/internal/model"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (name, password, role) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, user.Name, user.Password, user.Role).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByName(ctx context.Context, name string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, name, password, role, created_at, updated_at FROM users WHERE name = $1`

	err := r.db.GetContext(ctx, user, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by name: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, name, password, role, created_at, updated_at FROM users WHERE id = $1`

	err := r.db.GetContext(ctx, user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}
