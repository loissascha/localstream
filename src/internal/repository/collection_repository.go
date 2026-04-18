package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

var ErrCollectionNotFound = errors.New("collection not found")
var ErrCollectionMovieAlreadyExists = errors.New("collection movie already exists")
var ErrCollectionShowAlreadyExists = errors.New("collection show already exists")

type CollectionRepository interface {
	Create(ctx context.Context, collection *entity.Collection) error
	GetByIDForUser(ctx context.Context, id uuid.UUID, userID int64) (*entity.Collection, error)
	ListByUserID(ctx context.Context, userID int64) ([]entity.Collection, error)
	UpdateName(ctx context.Context, id uuid.UUID, userID int64, name string) error
	DeleteByIDForUser(ctx context.Context, id uuid.UUID, userID int64) error
	AddMovie(ctx context.Context, userID int64, collectionID, movieID uuid.UUID) error
	RemoveMovie(ctx context.Context, userID int64, collectionID, movieID uuid.UUID) error
	AddShow(ctx context.Context, userID int64, collectionID, showID uuid.UUID) error
	RemoveShow(ctx context.Context, userID int64, collectionID, showID uuid.UUID) error
	ListMovies(ctx context.Context, userID int64, collectionID uuid.UUID) ([]MovieSelectItem, error)
	ListShows(ctx context.Context, userID int64, collectionID uuid.UUID) ([]ShowSelectItem, error)
}
