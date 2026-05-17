package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type UserWatchstateRepository struct {
	db *sqlx.DB
}

func NewUserWatchstateRepository(db *sqlx.DB) *UserWatchstateRepository {
	return &UserWatchstateRepository{db: db}
}

func (r *UserWatchstateRepository) Upsert(ctx context.Context, watchstate *entity.UserWatchstate) error {
	const query = `
		INSERT INTO user_watchstates (user_id, show_id, season_id, episode_id, position, duration, finished)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, episode_id)
		DO UPDATE SET
			show_id = EXCLUDED.show_id,
			season_id = EXCLUDED.season_id,
			position = EXCLUDED.position,
			duration = EXCLUDED.duration,
			finished = EXCLUDED.finished,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowxContext(
		ctx,
		query,
		watchstate.UserID,
		watchstate.ShowID,
		watchstate.SeasonID,
		watchstate.EpisodeID,
		watchstate.Position,
		watchstate.Duration,
		watchstate.Finished,
	).Scan(&watchstate.ID, &watchstate.CreatedAt, &watchstate.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert user watchstate: %w", err)
	}

	return nil
}

func (r *UserWatchstateRepository) GetByUserAndEpisodeID(ctx context.Context, userId int64, episodeId uuid.UUID) (*entity.UserWatchstate, error) {
	const query = `
		SELECT id, user_id, show_id, season_id, episode_id, position, duration, finished, created_at, updated_at
		FROM user_watchstates
		WHERE user_id = $1 AND episode_id = $2
		LIMIT 1
	`

	var watchstate entity.UserWatchstate
	if err := r.db.GetContext(ctx, &watchstate, query, userId, episodeId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user watchstate by user and episode id: %w", err)
	}

	return &watchstate, nil
}

func (r *UserWatchstateRepository) ListByUserID(ctx context.Context, userId int64) ([]entity.UserWatchstate, error) {
	const query = `
		SELECT id, user_id, show_id, season_id, episode_id, position, duration, finished, created_at, updated_at
		FROM user_watchstates
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`

	var watchstates []entity.UserWatchstate
	if err := r.db.SelectContext(ctx, &watchstates, query, userId); err != nil {
		return nil, fmt.Errorf("list user watchstates by user id: %w", err)
	}

	return watchstates, nil
}

func (r *UserWatchstateRepository) ListLatestByShowForUserID(ctx context.Context, userId int64) ([]entity.UserWatchstate, error) {
	const query = `
		SELECT id, user_id, show_id, season_id, episode_id, position, duration, finished, created_at, updated_at
		FROM (
			SELECT DISTINCT ON (show_id) id, user_id, show_id, season_id, episode_id, position, duration, finished, created_at, updated_at
			FROM user_watchstates
			WHERE user_id = $1
			ORDER BY show_id, updated_at DESC
		) latest_watchstates
		ORDER BY updated_at DESC
		LIMIT 10
	`

	var watchstates []entity.UserWatchstate
	if err := r.db.SelectContext(ctx, &watchstates, query, userId); err != nil {
		return nil, fmt.Errorf("list latest user watchstates by show for user id: %w", err)
	}

	return watchstates, nil
}

func (r *UserWatchstateRepository) DeleteByUserAndEpisodeID(ctx context.Context, userId int64, episodeId uuid.UUID) error {
	const query = `
		DELETE FROM user_watchstates
		WHERE user_id = $1 AND episode_id = $2
	`

	if _, err := r.db.ExecContext(ctx, query, userId, episodeId); err != nil {
		return fmt.Errorf("delete user watchstate by user and episode id: %w", err)
	}

	return nil
}

var _ repository.UserWatchstateRepository = (*UserWatchstateRepository)(nil)
