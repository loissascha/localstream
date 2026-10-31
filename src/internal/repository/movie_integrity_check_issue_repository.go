package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

type MovieIntegrityCheckIssueRepository interface {
	Create(ctx context.Context, issue *entity.MovieIntegrityCheckIssue) error
	GetByMovieID(ctx context.Context, movieID uuid.UUID) ([]entity.MovieIntegrityCheckIssue, error)
	DeleteByMovieID(ctx context.Context, movieID uuid.UUID) error
}
