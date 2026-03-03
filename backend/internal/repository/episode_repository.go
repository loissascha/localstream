package repository

import (
	"context"

	"github.com/loissascha/localstream/internal/entity"
)

type EpisodeRepository interface {
	Create(ctx context.Context, episode *entity.Episode) error
	GetByPath(ctx context.Context, path string) (*entity.Episode, error)
}
