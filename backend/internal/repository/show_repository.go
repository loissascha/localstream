package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

type ShowRepository interface {
	Create(ctx context.Context, show *entity.Show) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Show, error)
	GetByPath(ctx context.Context, path string) (*entity.Show, error)
	List(ctx context.Context) ([]entity.Show, error)
}
