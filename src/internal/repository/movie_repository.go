package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

var ErrMovieNotFound = errors.New("movie not found")

type MovieRepository interface {
	Create(ctx context.Context, movie *entity.Movie) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Movie, error)
	GetByIDWithMetadata(ctx context.Context, id uuid.UUID, userID int64) (*MovieSelectItem, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	GetByPath(ctx context.Context, path string) (*entity.Movie, error)
	UpdateFetchSource(ctx context.Context, id uuid.UUID, fetchSource entity.FetchSource) error
	All(ctx context.Context) ([]entity.Movie, error)
	ListLatest(ctx context.Context, userID int64) ([]MovieSelectItem, error)
	List(ctx context.Context, userID int64) ([]MovieSelectItem, error)
	Search(ctx context.Context, query string, userID int64) ([]MovieSelectItem, error)
}

type MovieSelectItem struct {
	ID               uuid.UUID          `db:"id"`
	Name             string             `db:"name"`
	Year             int                `db:"year"`
	Description      string             `db:"description"`
	MediumImageUrl   string             `db:"medium_image_url"`
	BackdropImageUrl string             `db:"backdrop_image_url"`
	FetchSource      entity.FetchSource `db:"fetch_source"`
	Position         float64            `db:"position"`
	Duration         float64            `db:"duration"`
	Finished         bool               `db:"finished"`
	CreatedAt        time.Time          `db:"created_at"`
}
