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

func (r *EpisodeRepository) GetByID(ctx context.Context, episodeId uuid.UUID) (*repository.EpisodeWithMetadata, error) {
	const query = `
		SELECT e.id, coalesce(m.name, '') as "name", coalesce(m.summary, '') as "summary", coalesce(m.medium_image_url, '') as "medium_image_url", coalesce(m.original_image_url, '') as "original_image_url", coalesce(m.fetch_id, 0) as "fetch_id", e.season_id, e.number, e.path, e.created_at, e.fetch_source
		FROM episodes e
		LEFT JOIN episode_metadata m ON m.episode_id=e.id
		WHERE e.id = $1
		LIMIT 1
	`

	var episode repository.EpisodeWithMetadata
	if err := r.db.GetContext(ctx, &episode, query, episodeId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get episode by id: %w", err)
	}

	return &episode, nil
}

func (r *EpisodeRepository) DeleteByID(ctx context.Context, episodeId uuid.UUID) error {
	const query = `
		DELETE FROM episodes
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, episodeId)
	if err != nil {
		return fmt.Errorf("delete episode by id: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete episode by id rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrEpisodeNotFound
	}

	return nil
}

func (r *EpisodeRepository) GetByPathAndSeasonID(ctx context.Context, path string, seasonId uuid.UUID) (*entity.Episode, error) {
	const query = `
		SELECT *
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
		SELECT *
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

func (r *EpisodeRepository) ListBySeasonIDWithMetadata(ctx context.Context, seasonId uuid.UUID) ([]repository.EpisodeWithMetadata, error) {
	const query = `
		SELECT e.id, coalesce(m.name, '') as "name", coalesce(m.summary, '') as "summary", coalesce(m.medium_image_url, '') as "medium_image_url", coalesce(m.original_image_url, '') as "original_image_url", coalesce(m.fetch_id, 0) as "fetch_id", e.season_id, e.number, e.path, e.created_at, e.fetch_source
		FROM episodes e
		LEFT JOIN episode_metadata m ON m.episode_id=e.id
		WHERE e.season_id = $1
		ORDER BY e.number ASC
	`

	var episodes []repository.EpisodeWithMetadata
	if err := r.db.SelectContext(ctx, &episodes, query, seasonId); err != nil {
		return nil, fmt.Errorf("list episodes by season id: %w", err)
	}

	return episodes, nil
}

func (r *EpisodeRepository) GetBySeasonIDAndNumber(ctx context.Context, seasonId uuid.UUID, number int) (*repository.EpisodeWithMetadata, error) {
	const query = `
		SELECT e.id, coalesce(m.name, '') as "name", coalesce(m.summary, '') as "summary", coalesce(m.medium_image_url, '') as "medium_image_url", coalesce(m.original_image_url, '') as "original_image_url", coalesce(m.fetch_id, 0) as "fetch_id", e.season_id, e.number, e.path, e.created_at, e.fetch_source
		FROM episodes e
		LEFT JOIN episode_metadata m ON m.episode_id=e.id
		WHERE e.season_id = $1 AND e.number = $2
		LIMIT 1
	`

	var episode repository.EpisodeWithMetadata
	if err := r.db.GetContext(ctx, &episode, query, seasonId, number); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &episode, nil
}

func (r *EpisodeRepository) UpdateFetchSource(ctx context.Context, id uuid.UUID, fetchSource entity.FetchSource) error {
	if fetchSource == "" {
		fetchSource = entity.FetchSourceNone
	}

	const query = `
		UPDATE episodes
		SET fetch_source = $1
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, fetchSource, id)
	if err != nil {
		return fmt.Errorf("update episode fetch source: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update episode fetch source rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrEpisodeNotFound
	}

	return nil
}

var _ repository.EpisodeRepository = (*EpisodeRepository)(nil)
