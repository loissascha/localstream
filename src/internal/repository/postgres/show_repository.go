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

type ShowRepository struct {
	db *sqlx.DB
}

func NewShowRepository(db *sqlx.DB) *ShowRepository {
	return &ShowRepository{db: db}
}

func (r *ShowRepository) Create(ctx context.Context, show *entity.Show) error {
	fetchSource := show.FetchSource
	if fetchSource == "" {
		fetchSource = entity.FetchSourceNone
	}

	const query = `
		INSERT INTO shows (name, year, path, fetch_source)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRowxContext(ctx, query, show.Name, show.Year, show.Path, fetchSource).Scan(&show.ID, &show.CreatedAt)
	if err != nil {
		return fmt.Errorf("create show: %w", err)
	}

	show.FetchSource = fetchSource

	return nil
}

func (r *ShowRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Show, error) {
	const query = `
		SELECT id, name, year, path, created_at, fetch_source
		FROM shows
		WHERE id = $1
		LIMIT 1
	`

	var show entity.Show
	if err := r.db.GetContext(ctx, &show, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get show by id: %w", err)
	}

	return &show, nil
}

func (r *ShowRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM shows
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete show by id: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete show by id rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrShowNotFound
	}

	return nil
}

func (r *ShowRepository) GetByPath(ctx context.Context, path string) (*entity.Show, error) {
	const query = `
		SELECT id, name, year, path, created_at, fetch_source
		FROM shows
		WHERE path = $1
		LIMIT 1
	`

	var show entity.Show
	if err := r.db.GetContext(ctx, &show, query, path); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get show by path: %w", err)
	}

	return &show, nil
}

func (r *ShowRepository) UpdateFetchSource(ctx context.Context, id uuid.UUID, fetchSource entity.FetchSource) error {
	if fetchSource == "" {
		fetchSource = entity.FetchSourceNone
	}

	const query = `
		UPDATE shows
		SET fetch_source = $1
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, fetchSource, id)
	if err != nil {
		return fmt.Errorf("update show fetch source: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update show fetch source rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrShowNotFound
	}

	return nil
}
func (r *ShowRepository) ListLatest(ctx context.Context) ([]entity.Show, error) {
	const query = `
		SELECT *
		FROM shows
		ORDER BY created_at DESC
		LIMIT 10
	`

	var shows []entity.Show
	if err := r.db.SelectContext(ctx, &shows, query); err != nil {
		return nil, fmt.Errorf("list shows: %w", err)
	}

	return shows, nil
}

func (r *ShowRepository) List(ctx context.Context) ([]entity.Show, error) {
	const query = `
		SELECT *
		FROM shows
	`

	var shows []entity.Show
	if err := r.db.SelectContext(ctx, &shows, query); err != nil {
		return nil, fmt.Errorf("list shows: %w", err)
	}

	return shows, nil
}

func (r *ShowRepository) Search(ctx context.Context, query string) ([]entity.Show, error) {
	const stmt = `
		SELECT DISTINCT s.*
		FROM shows s
		WHERE s.name ILIKE $1
			OR EXISTS (
				SELECT 1
				FROM show_metadata sm
				WHERE sm.show_id = s.id
					AND sm.name ILIKE $1
			)
	`

	var shows []entity.Show
	if err := r.db.SelectContext(ctx, &shows, stmt, "%"+query+"%"); err != nil {
		return nil, fmt.Errorf("search shows: %w", err)
	}

	return shows, nil
}

var _ repository.ShowRepository = (*ShowRepository)(nil)
