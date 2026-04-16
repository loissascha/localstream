package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

var ErrShowMetadataNotFound = errors.New("show metadata not found")

type ShowMetadataRepository interface {
	Create(ctx context.Context, metadata *entity.ShowMetadata) error
	GetByShowID(ctx context.Context, showID uuid.UUID) ([]entity.ShowMetadata, error)
	DeleteOne(ctx context.Context, id uuid.UUID) error
}
