package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type MovieIntegrityCheckIssueRepository struct {
	db *sqlx.DB
}

func NewMovieIntegrityCheckIssueRepository(db *sqlx.DB) *MovieIntegrityCheckIssueRepository {
	return &MovieIntegrityCheckIssueRepository{db: db}
}

func (r *MovieIntegrityCheckIssueRepository) Create(ctx context.Context, issue *entity.MovieIntegrityCheckIssue) error {
	const query = `
		INSERT INTO movie_integrity_check_issues (movie_id, issue_level, issue)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	if err := r.db.QueryRowxContext(ctx, query, issue.MovieID, issue.IssueLevel, issue.Issue).Scan(&issue.ID); err != nil {
		return fmt.Errorf("create movie integrity check issue: %w", err)
	}
	return nil
}

func (r *MovieIntegrityCheckIssueRepository) GetByMovieIDs(ctx context.Context, movieIDs []uuid.UUID) ([]entity.MovieIntegrityCheckIssue, error) {
	const query = `
		SELECT *
		FROM movie_integrity_check_issues
		WHERE movie_id IN $1
		ORDER BY id ASC
	`

	issues := make([]entity.MovieIntegrityCheckIssue, 0)
	if err := r.db.GetContext(ctx, &issues, query, movieIDs); err != nil {
		return nil, fmt.Errorf("get movie integrity check issues by movie ids: %w", err)
	}
	return issues, nil
}

func (r *MovieIntegrityCheckIssueRepository) GetByMovieID(ctx context.Context, movieID uuid.UUID) ([]entity.MovieIntegrityCheckIssue, error) {
	const query = `
		SELECT id, movie_id, issue_level, issue
		FROM movie_integrity_check_issues
		WHERE movie_id = $1
		ORDER BY id ASC
	`

	issues := make([]entity.MovieIntegrityCheckIssue, 0)
	if err := r.db.SelectContext(ctx, &issues, query, movieID); err != nil {
		return nil, fmt.Errorf("get movie integrity check issues by movie id: %w", err)
	}
	return issues, nil
}

func (r *MovieIntegrityCheckIssueRepository) DeleteByMovieID(ctx context.Context, movieID uuid.UUID) error {
	const query = `
		DELETE FROM movie_integrity_check_issues
		WHERE movie_id = $1
	`

	if _, err := r.db.ExecContext(ctx, query, movieID); err != nil {
		return fmt.Errorf("delete movie integrity check issues by movie id: %w", err)
	}
	return nil
}

var _ repository.MovieIntegrityCheckIssueRepository = (*MovieIntegrityCheckIssueRepository)(nil)
