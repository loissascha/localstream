package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type MovieMetadataRepository struct {
	db *sqlx.DB
}

func NewMovieMetadataRepository(db *sqlx.DB) *MovieMetadataRepository {
	return &MovieMetadataRepository{db: db}
}

func (r *MovieMetadataRepository) Create(ctx context.Context, metadata *entity.MovieMetadata) error {
	const query = `
		INSERT INTO movie_metadata (movie_id, name, url, description, medium_image_url, backdrop_image_url, fetch_source)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.db.QueryRowxContext(
		ctx,
		query,
		metadata.MovieID,
		metadata.Name,
		metadata.Url,
		metadata.Description,
		metadata.MediumImageUrl,
		metadata.BackdropImageUrl,
		metadata.FetchSource,
	).Scan(&metadata.ID)
	if err != nil {
		return fmt.Errorf("create movie metadata: %w", err)
	}

	return nil
}

func (r *MovieMetadataRepository) GetByMovieID(ctx context.Context, movieID uuid.UUID) ([]entity.MovieMetadata, error) {
	const query = `
		SELECT *
		FROM movie_metadata
		WHERE movie_id = $1
	`

	var metadata []entity.MovieMetadata
	if err := r.db.SelectContext(ctx, &metadata, query, movieID); err != nil {
		return nil, fmt.Errorf("get movie metadata by movie id: %w", err)
	}

	return metadata, nil
}

func (r *MovieMetadataRepository) DeleteOne(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM movie_metadata
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
		return repository.ErrMovieMetadataNotFound
	}

	return nil
}

var _ repository.MovieMetadataRepository = (*MovieMetadataRepository)(nil)
