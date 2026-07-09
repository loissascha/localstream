package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

var ErrLibraryNotFound = errors.New("library not found")

type LibraryRepository interface {
	Create(ctx context.Context, library *entity.Library) error
	Update(ctx context.Context, library *entity.Library) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Library, error)
	List(ctx context.Context) ([]entity.Library, error)
}
