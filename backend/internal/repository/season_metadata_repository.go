package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

var ErrSeasonMetadataNotFound = errors.New("season metadata not found")

type SeasonMetadataRepository interface {
	Create(ctx context.Context, metadata *entity.SeasonMetadata) error
	GetByShowID(ctx context.Context, showID uuid.UUID) ([]entity.SeasonMetadata, error)
	DeleteOne(ctx context.Context, id uuid.UUID) error
}
