package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

var ErrMovieAudioStreamNotFound = errors.New("movie audio stream not found")

type MovieAudioStreamRepository interface {
	Create(ctx context.Context, stream *entity.MovieAudioStream) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.MovieAudioStream, error)
	ListByMovieID(ctx context.Context, movieID uuid.UUID) ([]entity.MovieAudioStream, error)
	Update(ctx context.Context, stream *entity.MovieAudioStream) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
