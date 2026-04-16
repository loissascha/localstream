package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

var ErrEpisodeNotFound = errors.New("episode not found")

type EpisodeRepository interface {
	Create(ctx context.Context, episode *entity.Episode) error
	GetByPathAndSeasonID(ctx context.Context, path string, seasonId uuid.UUID) (*entity.Episode, error)
	ListBySeasonID(ctx context.Context, seasonId uuid.UUID) ([]entity.Episode, error)
	GetByID(ctx context.Context, episodeId uuid.UUID) (*entity.Episode, error)
	DeleteByID(ctx context.Context, episodeId uuid.UUID) error
	GetBySeasonIDAndNumber(ctx context.Context, seasonId uuid.UUID, number int) (*entity.Episode, error)
	UpdateFetchSource(ctx context.Context, id uuid.UUID, fetchSource entity.FetchSource) error
}
