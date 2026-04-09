package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

var ErrMovieNotFound = errors.New("movie not found")

type MovieRepository interface {
	Create(ctx context.Context, movie *entity.Movie) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Movie, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	GetByPath(ctx context.Context, path string) (*entity.Movie, error)
	UpdateFetchSource(ctx context.Context, id uuid.UUID, fetchSource entity.FetchSource) error
	List(ctx context.Context) ([]entity.Movie, error)
}
