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
	GetByIDWithMetadata(ctx context.Context, id uuid.UUID) (*ShowSelectItem, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	GetByPath(ctx context.Context, path string) (*entity.Show, error)
	UpdateFetchSource(ctx context.Context, id uuid.UUID, fetchSource entity.FetchSource) error
	All(ctx context.Context) ([]entity.Show, error)
	List(ctx context.Context) ([]ShowSelectItem, error)
	ListLatest(ctx context.Context) ([]ShowSelectItem, error)
	Search(ctx context.Context, query string) ([]ShowSelectItem, error)
}

// TODO: SELECT s.id, COALESCE(m.name, s.name) as "name", s.fetch_source, s.path, COALESCE(m.description, ”) as "description", COALESCE(m.medium_image_url, ”) as "medium_image_url" from shows s LEFT JOIN show_metadata m ON m.show_id=s.id AND s.fetch_source!='multiple';
type ShowSelectItem struct {
	ID             uuid.UUID          `db:"id"`
	Name           string             `db:"name"`
	Year           int                `db:"year"`
	FetchSource    entity.FetchSource `db:"fetch_source"`
	Path           string             `db:"path"`
	Description    string             `db:"description"`
	MediumImageUrl string             `db:"medium_image_url"`
}
