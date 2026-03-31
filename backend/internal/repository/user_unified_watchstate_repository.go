package repository

import (
	"context"

	"github.com/loissascha/localstream/internal/entity"
)

type UserUnifiedWatchstateRepository interface {
	ListByUserID(ctx context.Context, userID int64) ([]entity.UserUnifiedWatchstate, error)
}
