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

type UserMovieWatchstateRepository struct {
	db *sqlx.DB
}

func NewUserMovieWatchstateRepository(db *sqlx.DB) *UserMovieWatchstateRepository {
	return &UserMovieWatchstateRepository{db: db}
}

func (r *UserMovieWatchstateRepository) Upsert(ctx context.Context, watchstate *entity.UserMovieWatchstate) error {
	const query = `
		INSERT INTO user_movie_watchstates (user_id, movie_id, position, duration, finished)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, movie_id)
		DO UPDATE SET
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
		watchstate.MovieID,
		watchstate.Position,
		watchstate.Duration,
		watchstate.Finished,
	).Scan(&watchstate.ID, &watchstate.CreatedAt, &watchstate.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert user movie watchstate: %w", err)
	}

	return nil
}

func (r *UserMovieWatchstateRepository) GetByUserAndMovieID(ctx context.Context, userID int64, movieID uuid.UUID) (*entity.UserMovieWatchstate, error) {
	const query = `
		SELECT id, user_id, movie_id, position, duration, finished, created_at, updated_at
		FROM user_movie_watchstates
		WHERE user_id = $1 AND movie_id = $2
		LIMIT 1
	`

	var watchstate entity.UserMovieWatchstate
	if err := r.db.GetContext(ctx, &watchstate, query, userID, movieID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user movie watchstate by user and movie id: %w", err)
	}

	return &watchstate, nil
}

func (r *UserMovieWatchstateRepository) ListByUserID(ctx context.Context, userID int64) ([]entity.UserMovieWatchstate, error) {
	const query = `
		SELECT id, user_id, movie_id, position, duration, finished, created_at, updated_at
		FROM user_movie_watchstates
		WHERE user_id = $1 AND finished = false
		ORDER BY updated_at DESC
		LIMIT 10
	`

	var watchstates []entity.UserMovieWatchstate
	if err := r.db.SelectContext(ctx, &watchstates, query, userID); err != nil {
		return nil, fmt.Errorf("list user movie watchstates by user id: %w", err)
	}

	return watchstates, nil
}

func (r *UserMovieWatchstateRepository) DeleteByUserAndMovieID(ctx context.Context, userID int64, movieID uuid.UUID) error {
	const query = `
		DELETE FROM user_movie_watchstates
		WHERE user_id = $1 AND movie_id = $2
	`

	if _, err := r.db.ExecContext(ctx, query, userID, movieID); err != nil {
		return fmt.Errorf("delete user movie watchstate by user and movie id: %w", err)
	}

	return nil
}

var _ repository.UserMovieWatchstateRepository = (*UserMovieWatchstateRepository)(nil)
