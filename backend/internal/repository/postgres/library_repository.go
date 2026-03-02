package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type LibraryRepository struct {
	db *sqlx.DB
}

func NewLibraryRepository(db *sqlx.DB) *LibraryRepository {
	return &LibraryRepository{db: db}
}

func (r *LibraryRepository) Create(ctx context.Context, library *entity.Library) error {
	const query = `
		INSERT INTO libraries (name, path, library_type)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRowxContext(ctx, query, library.Name, library.Path, library.LibraryType).Scan(&library.ID, &library.CreatedAt)
	if err != nil {
		return fmt.Errorf("create library: %w", err)
	}

	return nil
}

func (r *LibraryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Library, error) {
	const query = `
		SELECT id, name, path, library_type, created_at
		FROM libraries
		WHERE id = $1
	`

	var library entity.Library
	if err := r.db.GetContext(ctx, &library, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrLibraryNotFound
		}
		return nil, fmt.Errorf("get library by id: %w", err)
	}

	return &library, nil
}

func (r *LibraryRepository) List(ctx context.Context) ([]entity.Library, error) {
	const query = `
		SELECT id, name, path, library_type, created_at
		FROM libraries
		ORDER BY created_at ASC
	`

	var libraries []entity.Library
	if err := r.db.SelectContext(ctx, &libraries, query); err != nil {
		return nil, fmt.Errorf("list libraries: %w", err)
	}

	return libraries, nil
}

var _ repository.LibraryRepository = (*LibraryRepository)(nil)
