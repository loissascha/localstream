package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	const query = `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	err := r.db.QueryRowxContext(ctx, query, user.Username, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	const query = `
		SELECT id, username, password_hash, created_at
		FROM users
		WHERE id = $1
	`

	var user entity.User
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	const query = `
		SELECT id, username, password_hash, created_at
		FROM users
		WHERE username = $1
	`

	var user entity.User
	if err := r.db.GetContext(ctx, &user, query, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by username: %w", err)
	}

	return &user, nil
}

var _ repository.UserRepository = (*UserRepository)(nil)
