package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

var ErrEpisodeMetadataNotFound = errors.New("episode metadata not found")

type EpisodeMetadataRepository interface {
	Create(ctx context.Context, metadata *entity.EpisodeMetadata) error
	GetByShowID(ctx context.Context, showID uuid.UUID) ([]entity.EpisodeMetadata, error)
	GetByShowIDAndSeasonNumberAndEpisodeNumber(ctx context.Context, showID uuid.UUID, seasonNumber int, episodeNumber int) (*entity.EpisodeMetadata, error)
	DeleteOne(ctx context.Context, id uuid.UUID) error
}
