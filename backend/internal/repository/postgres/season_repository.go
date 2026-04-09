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

type SeasonRepository struct {
	db *sqlx.DB
}

func NewSeasonRepository(db *sqlx.DB) *SeasonRepository {
	return &SeasonRepository{db: db}
}

func (r *SeasonRepository) Create(ctx context.Context, season *entity.Season) error {
	fetchSource := season.FetchSource
	if fetchSource == "" {
		fetchSource = entity.FetchSourceNone
	}

	const query = `
		INSERT INTO seasons (show_id, number, path, fetch_source)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRowxContext(ctx, query, season.ShowID, season.Number, season.Path, fetchSource).Scan(&season.ID, &season.CreatedAt)
	if err != nil {
		return fmt.Errorf("create season: %w", err)
	}

	season.FetchSource = fetchSource

	return nil
}

func (r *SeasonRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Season, error) {
	const query = `
		SELECT id, show_id, number, path, created_at, fetch_source
		FROM seasons
		WHERE id = $1
		LIMIT 1
	`
	var season entity.Season
	if err := r.db.GetContext(ctx, &season, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get season by path: %w", err)
	}

	return &season, nil
}

func (r *SeasonRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM seasons
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete season by id: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete season by id rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrSeasonNotFound
	}

	return nil
}

func (r *SeasonRepository) GetByPathAndShowID(ctx context.Context, path string, showId uuid.UUID) (*entity.Season, error) {
	const query = `
		SELECT id, show_id, number, path, created_at, fetch_source
		FROM seasons
		WHERE path = $1 AND show_id = $2
		LIMIT 1
	`

	var season entity.Season
	if err := r.db.GetContext(ctx, &season, query, path, showId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get season by path: %w", err)
	}

	return &season, nil
}

func (r *SeasonRepository) ListByShowID(ctx context.Context, showId uuid.UUID) ([]entity.Season, error) {
	const query = `
		SELECT id, show_id, number, path, created_at, fetch_source
		FROM seasons
		WHERE show_id = $1
		ORDER BY number ASC
	`

	var seasons []entity.Season
	if err := r.db.SelectContext(ctx, &seasons, query, showId); err != nil {
		return nil, fmt.Errorf("list seasons by show id: %w", err)
	}

	return seasons, nil
}

var _ repository.SeasonRepository = (*SeasonRepository)(nil)
