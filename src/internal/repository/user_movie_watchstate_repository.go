package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

type UserMovieWatchstateRepository interface {
	Upsert(ctx context.Context, watchstate *entity.UserMovieWatchstate) error
	GetByUserAndMovieID(ctx context.Context, userID int64, movieID uuid.UUID) (*entity.UserMovieWatchstate, error)
	ListByUserID(ctx context.Context, userID int64) ([]entity.UserMovieWatchstate, error)
	DeleteByUserAndMovieID(ctx context.Context, userID int64, movieID uuid.UUID) error
}
