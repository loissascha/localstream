package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

type EpisodeRepository interface {
	Create(ctx context.Context, episode *entity.Episode) error
	GetByPathAndSeasonID(ctx context.Context, path string, seasonId uuid.UUID) (*entity.Episode, error)
}
