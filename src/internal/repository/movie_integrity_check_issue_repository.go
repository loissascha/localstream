package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

type MovieIntegrityCheckIssueRepository interface {
	Create(ctx context.Context, show *entity.MovieIntegrityCheckIssue) error
	GetByMovieID(ctx context.Context, movieId uuid.UUID) (*entity.MovieIntegrityCheckIssue, error)
	DeleteByMovieID(ctx context.Context, movieId uuid.UUID) error
}
