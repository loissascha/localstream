package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type UserUnifiedWatchstateRepository struct {
	db *sqlx.DB
}

func NewUserUnifiedWatchstateRepository(db *sqlx.DB) *UserUnifiedWatchstateRepository {
	return &UserUnifiedWatchstateRepository{db: db}
}

func (r *UserUnifiedWatchstateRepository) ListByUserID(ctx context.Context, userID int64) ([]entity.UserUnifiedWatchstate, error) {
	const query = `
		SELECT id, user_id, media_type, media_id, position, duration, finished, created_at, updated_at
		FROM user_all_watchstates
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`

	var watchstates []entity.UserUnifiedWatchstate
	if err := r.db.SelectContext(ctx, &watchstates, query, userID); err != nil {
		return nil, fmt.Errorf("list unified watchstates by user id: %w", err)
	}

	return watchstates, nil
}

var _ repository.UserUnifiedWatchstateRepository = (*UserUnifiedWatchstateRepository)(nil)
