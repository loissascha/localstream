package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

type ShowMetadataRepository interface {
	Create(ctx context.Context, metadata *entity.ShowMetadata) error
	GetByShowID(ctx context.Context, showID uuid.UUID) ([]entity.ShowMetadata, error)
}
