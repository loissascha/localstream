package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

var ErrMovieMetadataNotFound = errors.New("movie metadata not found")

type MovieMetadataRepository interface {
	Create(ctx context.Context, metadata *entity.MovieMetadata) error
	GetByMovieID(ctx context.Context, movieID uuid.UUID) ([]entity.MovieMetadata, error)
	DeleteOne(ctx context.Context, id uuid.UUID) error
}
