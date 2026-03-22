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

type EpisodeRepository struct {
	db *sqlx.DB
}

func NewEpisodeRepository(db *sqlx.DB) *EpisodeRepository {
	return &EpisodeRepository{db: db}
}

func (r *EpisodeRepository) Create(ctx context.Context, episode *entity.Episode) error {
	fetchSource := episode.FetchSource
	if fetchSource == "" {
		fetchSource = entity.FetchSourceNone
	}

	const query = `
		INSERT INTO episodes (id, season_id, number, path, fetch_source)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	err = r.db.QueryRowxContext(ctx, query, id, episode.SeasonID, episode.Number, episode.Path, fetchSource).Scan(&episode.ID, &episode.CreatedAt)
	if err != nil {
		return fmt.Errorf("create episode: %w", err)
	}

	episode.FetchSource = fetchSource

	return nil
}

func (r *EpisodeRepository) GetByID(ctx context.Context, episodeId uuid.UUID) (*entity.Episode, error) {
	const query = `
		SELECT id, season_id, number, path, created_at, fetch_source
		FROM episodes
		WHERE id = $1
		LIMIT 1
	`

	var episode entity.Episode
	if err := r.db.GetContext(ctx, &episode, query, episodeId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get episode by id: %w", err)
	}

	return &episode, nil
}

func (r *EpisodeRepository) GetByPathAndSeasonID(ctx context.Context, path string, seasonId uuid.UUID) (*entity.Episode, error) {
	const query = `
		SELECT id, season_id, number, path, created_at, fetch_source
		FROM episodes
		WHERE path = $1 AND season_id = $2
		LIMIT 1
	`

	var episode entity.Episode
	if err := r.db.GetContext(ctx, &episode, query, path, seasonId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get episode by path: %w", err)
	}

	return &episode, nil
}

func (r *EpisodeRepository) ListBySeasonID(ctx context.Context, seasonId uuid.UUID) ([]entity.Episode, error) {
	const query = `
		SELECT id, season_id, number, path, created_at, fetch_source
		FROM episodes
		WHERE season_id = $1
		ORDER BY number ASC
	`

	var episodes []entity.Episode
	if err := r.db.SelectContext(ctx, &episodes, query, seasonId); err != nil {
		return nil, fmt.Errorf("list episodes by season id: %w", err)
	}

	return episodes, nil
}

func (r *EpisodeRepository) GetBySeasonIDAndNumber(ctx context.Context, seasonId uuid.UUID, number int) (*entity.Episode, error) {
	const query = `
		SELECT id, season_id, number, path, created_at, fetch_source
		FROM episodes
		WHERE season_id = $1 AND number = $2
		LIMIT 1
	`

	var episode entity.Episode
	if err := r.db.GetContext(ctx, &episode, query, seasonId, number); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &episode, nil
}

var _ repository.EpisodeRepository = (*EpisodeRepository)(nil)
