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

type EpisodeMetadataRepository struct {
	db *sqlx.DB
}

func NewEpisodeMetadataRepository(db *sqlx.DB) *EpisodeMetadataRepository {
	return &EpisodeMetadataRepository{db: db}
}

func (r *EpisodeMetadataRepository) Create(ctx context.Context, metadata *entity.EpisodeMetadata) error {
	const query = `
		INSERT INTO episode_metadata (episode_id, url, name, number, summary, medium_image_url, original_image_url, fetch_id, fetch_source)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	err := r.db.QueryRowxContext(
		ctx,
		query,
		metadata.EpisodeID,
		metadata.Url,
		metadata.Name,
		metadata.Number,
		metadata.Summary,
		metadata.MediumImageUrl,
		metadata.OriginalImageUrl,
		metadata.FetchID,
		metadata.FetchSource,
	).Scan(&metadata.ID)
	if err != nil {
		return fmt.Errorf("create episode metadata: %w", err)
	}

	return nil
}

func (r *EpisodeMetadataRepository) GetByShowID(ctx context.Context, showID uuid.UUID) ([]entity.EpisodeMetadata, error) {
	const query = `
		SELECT em.*
		FROM episode_metadata em
		INNER JOIN episodes e ON e.id = em.episode_id
		INNER JOIN seasons s ON s.id = e.season_id
		WHERE s.show_id = $1
		ORDER BY s.number ASC, e.number ASC
	`

	var metadata []entity.EpisodeMetadata
	if err := r.db.SelectContext(ctx, &metadata, query, showID); err != nil {
		return nil, fmt.Errorf("get episode metadata by show id: %w", err)
	}

	return metadata, nil
}

func (r *EpisodeMetadataRepository) GetByShowIDAndSeasonNumberAndEpisodeNumber(ctx context.Context, showID uuid.UUID, seasonNumber int, episodeNumber int) (*entity.EpisodeMetadata, error) {
	const query = `
		SELECT em.*
		FROM episode_metadata em
		INNER JOIN episodes e ON e.id = em.episode_id
		INNER JOIN seasons s ON s.id = e.season_id
		WHERE s.show_id = $1 AND s.number = $2 AND e.number = $3
		ORDER BY em.id ASC
		LIMIT 1
	`

	var metadata entity.EpisodeMetadata
	if err := r.db.GetContext(ctx, &metadata, query, showID, seasonNumber, episodeNumber); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get episode metadata by show id, season number and episode number: %w", err)
	}

	return &metadata, nil
}

func (r *EpisodeMetadataRepository) DeleteOne(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM episode_metadata
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repository.ErrEpisodeMetadataNotFound
	}

	return nil
}

var _ repository.EpisodeMetadataRepository = (*EpisodeMetadataRepository)(nil)
