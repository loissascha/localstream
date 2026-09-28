package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

var ErrMovieVideoStreamNotFound = errors.New("movie video stream not found")

type MovieVideoStreamRepository interface {
	Create(ctx context.Context, stream *entity.MovieVideoStream) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.MovieVideoStream, error)
	ListByMovieID(ctx context.Context, movieID uuid.UUID) ([]entity.MovieVideoStream, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
