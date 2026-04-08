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

type MovieRepository struct {
	db *sqlx.DB
}

func NewMovieRepository(db *sqlx.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Movie, error) {
	const query = `
		SELECT *
		FROM movies
		WHERE id = $1
		LIMIT 1
	`

	var movie entity.Movie
	if err := r.db.GetContext(ctx, &movie, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get movie by id: %w", err)
	}

	return &movie, nil
}

func (r *MovieRepository) GetByPath(ctx context.Context, path string) (*entity.Movie, error) {
	const query = `
		SELECT *
		FROM movies
		WHERE path = $1
		LIMIT 1
	`

	var movie entity.Movie
	if err := r.db.GetContext(ctx, &movie, query, path); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get movie by path: %w", err)
	}

	return &movie, nil
}

func (r *MovieRepository) List(ctx context.Context) ([]entity.Movie, error) {
	const query = `
		SELECT *
		FROM movies
	`

	var movies []entity.Movie
	if err := r.db.SelectContext(ctx, &movies, query); err != nil {
		return nil, fmt.Errorf("list movies: %w", err)
	}

	return movies, nil
}

func (r *MovieRepository) UpdateFetchSource(ctx context.Context, id uuid.UUID, fetchSource entity.FetchSource) error {
	if fetchSource == "" {
		fetchSource = entity.FetchSourceNone
	}

	const query = `
		UPDATE movies
		SET fetch_source = $1
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, fetchSource, id)
	if err != nil {
		return fmt.Errorf("update movie fetch source: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update movie fetch source rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrMovieNotFound
	}

	return nil
}

func (r *MovieRepository) Create(ctx context.Context, movie *entity.Movie) error {
	fetchSource := movie.FetchSource
	if fetchSource == "" {
		fetchSource = entity.FetchSourceNone
	}
	movie.FetchSource = fetchSource

	const query = `
		INSERT INTO movies (id, name, year, description, path, created_at, fetch_source)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	movie.ID = id

	_, err = r.db.ExecContext(ctx, query, movie.ID, movie.Name, movie.Year, movie.Description, movie.Path, movie.CreatedAt, movie.FetchSource)
	if err != nil {
		return fmt.Errorf("create movie: %w", err)
	}

	return nil
}

var _ repository.MovieRepository = (*MovieRepository)(nil)
