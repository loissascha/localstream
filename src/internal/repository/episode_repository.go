package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

var ErrEpisodeNotFound = errors.New("episode not found")

type EpisodeRepository interface {
	Create(ctx context.Context, episode *entity.Episode) error
	GetByPathAndSeasonID(ctx context.Context, path string, seasonId uuid.UUID) (*entity.Episode, error)
	ListBySeasonID(ctx context.Context, seasonId uuid.UUID) ([]entity.Episode, error)
	ListBySeasonIDWithMetadata(ctx context.Context, seasonId uuid.UUID) ([]EpisodeWithMetadata, error)
	GetByID(ctx context.Context, episodeId uuid.UUID) (*EpisodeWithMetadata, error)
	DeleteByID(ctx context.Context, episodeId uuid.UUID) error
	GetBySeasonIDAndNumber(ctx context.Context, seasonId uuid.UUID, number int) (*EpisodeWithMetadata, error)
	UpdateFetchSource(ctx context.Context, id uuid.UUID, fetchSource entity.FetchSource) error
}

type EpisodeWithMetadata struct {
	ID               uuid.UUID          `db:"id"`
	SeasonID         uuid.UUID          `db:"season_id"`
	Number           int                `db:"number"`
	Path             string             `db:"path"`
	CreatedAt        time.Time          `db:"created_at"`
	FetchSource      entity.FetchSource `db:"fetch_source"`
	Name             string             `db:"name"`
	Summary          string             `db:"summary"`
	MediumImageUrl   string             `db:"medium_image_url"`
	OriginalImageUrl string             `db:"original_image_url"`
	FetchID          int                `db:"fetch_id"`
}
