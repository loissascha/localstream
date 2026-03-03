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
		INSERT INTO seasons (show_id, name, path, fetch_source)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRowxContext(ctx, query, season.ShowID, season.Name, season.Path, fetchSource).Scan(&season.ID, &season.CreatedAt)
	if err != nil {
		return fmt.Errorf("create season: %w", err)
	}

	season.FetchSource = fetchSource

	return nil
}

func (r *SeasonRepository) GetByPath(ctx context.Context, path string) (*entity.Season, error) {
	const query = `
		SELECT id, show_id, name, path, created_at, fetch_source
		FROM seasons
		WHERE path = $1
		LIMIT 1
	`

	var season entity.Season
	if err := r.db.GetContext(ctx, &season, query, path); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &entity.Season{}, nil
		}
		return nil, fmt.Errorf("get season by path: %w", err)
	}

	return &season, nil
}

var _ repository.SeasonRepository = (*SeasonRepository)(nil)
