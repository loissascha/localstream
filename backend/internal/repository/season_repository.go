package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

type SeasonRepository interface {
	Create(ctx context.Context, season *entity.Season) error
	GetByPathAndShowID(ctx context.Context, path string, showId uuid.UUID) (*entity.Season, error)
}
