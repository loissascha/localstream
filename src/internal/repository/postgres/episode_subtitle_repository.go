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

type EpisodeSubtitleRepository struct {
	db *sqlx.DB
}

func NewEpisodeSubtitleRepository(db *sqlx.DB) *EpisodeSubtitleRepository {
	return &EpisodeSubtitleRepository{db: db}
}

func (r *EpisodeSubtitleRepository) Create(ctx context.Context, subtitle *entity.EpisodeSubtitle) error {
	const query = `
		INSERT INTO episode_subtitles (id, episode_id, path, name, lang_short, lang)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	subtitle.ID = id

	_, err = r.db.ExecContext(ctx, query, subtitle.ID, subtitle.EpisodeID, subtitle.Path, subtitle.Name, subtitle.LangShort, subtitle.Lang)
	if err != nil {
		return fmt.Errorf("create episode subtitle: %w", err)
	}

	return nil
}

func (r *EpisodeSubtitleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.EpisodeSubtitle, error) {
	const query = `
		SELECT *
		FROM episode_subtitles
		WHERE id = $1
		LIMIT 1
	`

	var subtitle entity.EpisodeSubtitle
	if err := r.db.GetContext(ctx, &subtitle, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get episode subtitle by id: %w", err)
	}

	return &subtitle, nil
}

func (r *EpisodeSubtitleRepository) GetByPath(ctx context.Context, path string) (*entity.EpisodeSubtitle, error) {
	const query = `
		SELECT *
		FROM episode_subtitles
		WHERE path = $1
		LIMIT 1
	`

	var subtitle entity.EpisodeSubtitle
	if err := r.db.GetContext(ctx, &subtitle, query, path); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get episode subtitle by path: %w", err)
	}

	return &subtitle, nil
}

func (r *EpisodeSubtitleRepository) ListByEpisodeID(ctx context.Context, episodeID uuid.UUID) ([]entity.EpisodeSubtitle, error) {
	const query = `
		SELECT *
		FROM episode_subtitles
		WHERE episode_id = $1
		ORDER BY name ASC
	`

	var subtitles []entity.EpisodeSubtitle
	if err := r.db.SelectContext(ctx, &subtitles, query, episodeID); err != nil {
		return nil, fmt.Errorf("list episode subtitles by episode id: %w", err)
	}

	return subtitles, nil
}

func (r *EpisodeSubtitleRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM episode_subtitles
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete episode subtitle by id: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete episode subtitle by id rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrEpisodeSubtitleNotFound
	}

	return nil
}

var _ repository.EpisodeSubtitleRepository = (*EpisodeSubtitleRepository)(nil)
