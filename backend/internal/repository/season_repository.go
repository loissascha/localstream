package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

var ErrSeasonNotFound = errors.New("season not found")

type SeasonRepository interface {
	Create(ctx context.Context, season *entity.Season) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Season, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	GetByPathAndShowID(ctx context.Context, path string, showId uuid.UUID) (*entity.Season, error)
	ListByShowID(ctx context.Context, showId uuid.UUID) ([]entity.Season, error)
}
