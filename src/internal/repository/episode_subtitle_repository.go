package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

var ErrEpisodeSubtitleNotFound = errors.New("episode subtitle not found")

type EpisodeSubtitleRepository interface {
	Create(ctx context.Context, subtitle *entity.EpisodeSubtitle) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.EpisodeSubtitle, error)
	GetByPath(ctx context.Context, path string) (*entity.EpisodeSubtitle, error)
	ListByEpisodeID(ctx context.Context, episodeID uuid.UUID) ([]entity.EpisodeSubtitle, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
