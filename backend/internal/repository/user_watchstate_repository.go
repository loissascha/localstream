package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
)

type UserWatchstateRepository interface {
	Upsert(ctx context.Context, watchstate *entity.UserWatchstate) error
	GetByUserAndEpisodeID(ctx context.Context, userId int64, episodeId uuid.UUID) (*entity.UserWatchstate, error)
	ListByUserID(ctx context.Context, userId int64) ([]entity.UserWatchstate, error)
	ListLatestByShowForUserID(ctx context.Context, userId int64) ([]entity.UserWatchstate, error)
	DeleteByUserAndEpisodeID(ctx context.Context, userId int64, episodeId uuid.UUID) error
}
