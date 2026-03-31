package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

type MovieRepository interface {
	Create(ctx context.Context, movie *entity.Movie) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Movie, error)
	GetByPath(ctx context.Context, path string) (*entity.Movie, error)
	List(ctx context.Context) ([]entity.Movie, error)
}
