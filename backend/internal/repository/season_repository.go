package repository

import (
	"context"

	"github.com/loissascha/localstream/internal/entity"
)

type SeasonRepository interface {
	Create(ctx context.Context, season *entity.Season) error
	GetByPath(ctx context.Context, path string) (*entity.Season, error)
}
