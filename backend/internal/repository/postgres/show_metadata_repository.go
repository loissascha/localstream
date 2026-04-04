package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type ShowMetadataRepository struct {
	db *sqlx.DB
}

func NewShowMetadataRepository(db *sqlx.DB) *ShowMetadataRepository {
	return &ShowMetadataRepository{db: db}
}

func (r *ShowMetadataRepository) Create(ctx context.Context, metadata *entity.ShowMetadata) error {
	const query = `
		INSERT INTO show_metadata (show_id, name, url, description, medium_image_url, original_image_url, fetch_source)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.db.QueryRowxContext(
		ctx,
		query,
		metadata.ShowID,
		metadata.Name,
		metadata.Url,
		metadata.Description,
		metadata.MediumImageUrl,
		metadata.OriginalImageUrl,
		metadata.FetchSource,
	).Scan(&metadata.ID)
	if err != nil {
		return fmt.Errorf("create show metadata: %w", err)
	}

	return nil
}

func (r *ShowMetadataRepository) GetByShowID(ctx context.Context, showID uuid.UUID) ([]entity.ShowMetadata, error) {
	const query = `
		SELECT *
		FROM show_metadata
		WHERE show_id = $1
	`

	var metadata []entity.ShowMetadata
	if err := r.db.SelectContext(ctx, &metadata, query, showID); err != nil {
		return nil, fmt.Errorf("get show metadata by show id: %w", err)
	}

	return metadata, nil
}

var _ repository.ShowMetadataRepository = (*ShowMetadataRepository)(nil)
