package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

var ErrMovieSubtitleNotFound = errors.New("movie subtitle not found")

type MovieSubtitleRepository interface {
	Create(ctx context.Context, subtitle *entity.MovieSubtitle) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.MovieSubtitle, error)
	GetByPath(ctx context.Context, path string) (*entity.MovieSubtitle, error)
	ListByMovieID(ctx context.Context, movieID uuid.UUID) ([]entity.MovieSubtitle, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
