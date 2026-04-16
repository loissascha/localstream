package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

var ErrShowNotFound = errors.New("show not found")

type ShowRepository interface {
	Create(ctx context.Context, show *entity.Show) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Show, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	GetByPath(ctx context.Context, path string) (*entity.Show, error)
	UpdateFetchSource(ctx context.Context, id uuid.UUID, fetchSource entity.FetchSource) error
	List(ctx context.Context) ([]entity.Show, error)
	ListLatest(ctx context.Context) ([]entity.Show, error)
	Search(ctx context.Context, query string) ([]entity.Show, error)
}
