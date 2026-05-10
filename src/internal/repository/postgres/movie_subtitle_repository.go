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

type MovieSubtitleRepository struct {
	db *sqlx.DB
}

func NewMovieSubtitleRepository(db *sqlx.DB) *MovieSubtitleRepository {
	return &MovieSubtitleRepository{db: db}
}

func (r *MovieSubtitleRepository) Create(ctx context.Context, subtitle *entity.MovieSubtitle) error {
	const query = `
		INSERT INTO movie_subtitles (id, movie_id, path, name, lang_short, lang)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	subtitle.ID = id

	_, err = r.db.ExecContext(ctx, query, subtitle.ID, subtitle.MovieID, subtitle.Path, subtitle.Name, subtitle.LangShort, subtitle.Lang)
	if err != nil {
		return fmt.Errorf("create movie subtitle: %w", err)
	}

	return nil
}

func (r *MovieSubtitleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.MovieSubtitle, error) {
	const query = `
		SELECT *
		FROM movie_subtitles
		WHERE id = $1
		LIMIT 1
	`

	var subtitle entity.MovieSubtitle
	if err := r.db.GetContext(ctx, &subtitle, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get movie subtitle by id: %w", err)
	}

	return &subtitle, nil
}

func (r *MovieSubtitleRepository) GetByPath(ctx context.Context, path string) (*entity.MovieSubtitle, error) {
	const query = `
		SELECT *
		FROM movie_subtitles
		WHERE path = $1
		LIMIT 1
	`

	var subtitle entity.MovieSubtitle
	if err := r.db.GetContext(ctx, &subtitle, query, path); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get movie subtitle by path: %w", err)
	}

	return &subtitle, nil
}

func (r *MovieSubtitleRepository) ListByMovieID(ctx context.Context, movieID uuid.UUID) ([]entity.MovieSubtitle, error) {
	const query = `
		SELECT *
		FROM movie_subtitles
		WHERE movie_id = $1
		ORDER BY name ASC
	`

	var subtitles []entity.MovieSubtitle
	if err := r.db.SelectContext(ctx, &subtitles, query, movieID); err != nil {
		return nil, fmt.Errorf("list movie subtitles by movie id: %w", err)
	}

	return subtitles, nil
}

func (r *MovieSubtitleRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM movie_subtitles
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete movie subtitle by id: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete movie subtitle by id rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrMovieSubtitleNotFound
	}

	return nil
}

var _ repository.MovieSubtitleRepository = (*MovieSubtitleRepository)(nil)
