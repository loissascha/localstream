package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type SeasonMetadataRepository struct {
	db *sqlx.DB
}

func NewSeasonMetadataRepository(db *sqlx.DB) *SeasonMetadataRepository {
	return &SeasonMetadataRepository{db: db}
}

func (r *SeasonMetadataRepository) Create(ctx context.Context, metadata *entity.SeasonMetadata) error {
	const query = `
		INSERT INTO season_metadata (show_id, show_metadata_id, url, number, summary, premiere_date, medium_image_url, original_image_url, fetch_id, fetch_source)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	err := r.db.QueryRowxContext(
		ctx,
		query,
		metadata.ShowID,
		metadata.ShowMetadataID,
		metadata.Url,
		metadata.Number,
		metadata.Summary,
		metadata.PremiereDate,
		metadata.MediumImageUrl,
		metadata.OriginalImageUrl,
		metadata.FetchID,
		metadata.FetchSource,
	).Scan(&metadata.ID)
	if err != nil {
		return fmt.Errorf("create season metadata: %w", err)
	}

	return nil
}

func (r *SeasonMetadataRepository) GetByShowID(ctx context.Context, showID uuid.UUID) ([]entity.SeasonMetadata, error) {
	const query = `
		SELECT *
		FROM season_metadata
		WHERE show_id = $1
		ORDER BY number ASC
	`

	var metadata []entity.SeasonMetadata
	if err := r.db.SelectContext(ctx, &metadata, query, showID); err != nil {
		return nil, fmt.Errorf("get season metadata by show id: %w", err)
	}

	return metadata, nil
}

func (r *SeasonMetadataRepository) DeleteOne(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM season_metadata
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
		return repository.ErrSeasonMetadataNotFound
	}

	return nil
}

var _ repository.SeasonMetadataRepository = (*SeasonMetadataRepository)(nil)
