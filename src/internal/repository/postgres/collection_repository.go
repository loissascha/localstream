package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type CollectionRepository struct {
	db *sqlx.DB
}

func NewCollectionRepository(db *sqlx.DB) *CollectionRepository {
	return &CollectionRepository{db: db}
}

func (r *CollectionRepository) Create(ctx context.Context, collection *entity.Collection) error {
	const query = `
		INSERT INTO collections (user_id, name)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowxContext(ctx, query, collection.UserID, collection.Name).Scan(
		&collection.ID,
		&collection.CreatedAt,
		&collection.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create collection: %w", err)
	}

	return nil
}

func (r *CollectionRepository) GetByIDForUser(ctx context.Context, id uuid.UUID, userID int64) (*entity.Collection, error) {
	const query = `
		SELECT id, user_id, name, created_at, updated_at
		FROM collections
		WHERE id = $1 AND user_id = $2
		LIMIT 1
	`

	var collection entity.Collection
	if err := r.db.GetContext(ctx, &collection, query, id, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get collection by id for user: %w", err)
	}

	return &collection, nil
}

func (r *CollectionRepository) ListByUserID(ctx context.Context, userID int64) ([]entity.Collection, error) {
	const query = `
		SELECT id, user_id, name, created_at, updated_at
		FROM collections
		WHERE user_id = $1
		ORDER BY updated_at DESC, created_at DESC
	`

	var collections []entity.Collection
	if err := r.db.SelectContext(ctx, &collections, query, userID); err != nil {
		return nil, fmt.Errorf("list collections by user id: %w", err)
	}

	return collections, nil
}

func (r *CollectionRepository) UpdateName(ctx context.Context, id uuid.UUID, userID int64, name string) error {
	const query = `
		UPDATE collections
		SET name = $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
	`

	result, err := r.db.ExecContext(ctx, query, name, id, userID)
	if err != nil {
		return fmt.Errorf("update collection name: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update collection name rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrCollectionNotFound
	}

	return nil
}

func (r *CollectionRepository) DeleteByIDForUser(ctx context.Context, id uuid.UUID, userID int64) error {
	const query = `
		DELETE FROM collections
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete collection by id for user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete collection by id for user rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrCollectionNotFound
	}

	return nil
}

func (r *CollectionRepository) AddMovie(ctx context.Context, userID int64, collectionID, movieID uuid.UUID) error {
	const query = `
		INSERT INTO collection_movies (collection_id, movie_id)
		SELECT c.id, $3
		FROM collections c
		WHERE c.id = $2 AND c.user_id = $1
	`

	result, err := r.db.ExecContext(ctx, query, userID, collectionID, movieID)
	if err != nil {
		if isUniqueViolation(err) {
			return repository.ErrCollectionMovieAlreadyExists
		}
		return fmt.Errorf("add movie to collection: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("add movie to collection rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrCollectionNotFound
	}

	return nil
}

func (r *CollectionRepository) RemoveMovie(ctx context.Context, userID int64, collectionID, movieID uuid.UUID) error {
	const query = `
		DELETE FROM collection_movies cm
		USING collections c
		WHERE cm.collection_id = c.id
			AND c.user_id = $1
			AND cm.collection_id = $2
			AND cm.movie_id = $3
	`

	result, err := r.db.ExecContext(ctx, query, userID, collectionID, movieID)
	if err != nil {
		return fmt.Errorf("remove movie from collection: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("remove movie from collection rows affected: %w", err)
	}

	if rowsAffected > 0 {
		return nil
	}

	return r.ensureCollectionExistsForUser(ctx, collectionID, userID)
}

func (r *CollectionRepository) AddShow(ctx context.Context, userID int64, collectionID, showID uuid.UUID) error {
	const query = `
		INSERT INTO collection_shows (collection_id, show_id)
		SELECT c.id, $3
		FROM collections c
		WHERE c.id = $2 AND c.user_id = $1
	`

	result, err := r.db.ExecContext(ctx, query, userID, collectionID, showID)
	if err != nil {
		if isUniqueViolation(err) {
			return repository.ErrCollectionShowAlreadyExists
		}
		return fmt.Errorf("add show to collection: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("add show to collection rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrCollectionNotFound
	}

	return nil
}

func (r *CollectionRepository) RemoveShow(ctx context.Context, userID int64, collectionID, showID uuid.UUID) error {
	const query = `
		DELETE FROM collection_shows cs
		USING collections c
		WHERE cs.collection_id = c.id
			AND c.user_id = $1
			AND cs.collection_id = $2
			AND cs.show_id = $3
	`

	result, err := r.db.ExecContext(ctx, query, userID, collectionID, showID)
	if err != nil {
		return fmt.Errorf("remove show from collection: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("remove show from collection rows affected: %w", err)
	}

	if rowsAffected > 0 {
		return nil
	}

	return r.ensureCollectionExistsForUser(ctx, collectionID, userID)
}

func (r *CollectionRepository) ListMovies(ctx context.Context, userID int64, collectionID uuid.UUID) ([]repository.MovieSelectItem, error) {
	const query = `
		SELECT
			m.id,
			COALESCE(umw.position, 0) as "position",
			COALESCE(umw.duration, 0) as "duration",
			COALESCE(umw.finished, false) as "finished",
			COALESCE(mm.name, m.name) as "name",
			COALESCE(mm.release_year, m.year) as "year",
			COALESCE(mm.description, m.description) as "description",
			COALESCE(mm.medium_image_url, '') as "medium_image_url",
			COALESCE(mm.backdrop_image_url, '') as "backdrop_image_url",
			m.fetch_source
		FROM collection_movies cm
		JOIN collections c
			ON c.id = cm.collection_id
		JOIN movies m
			ON m.id = cm.movie_id
		LEFT JOIN movie_metadata mm
			ON mm.movie_id = m.id AND m.fetch_source != 'multiple'
		LEFT JOIN user_movie_watchstates umw
			ON umw.movie_id = m.id AND umw.user_id = $1
		WHERE c.user_id = $1 AND c.id = $2
		ORDER BY COALESCE(mm.release_year, m.year) DESC, LOWER(COALESCE(mm.name, m.name)) ASC
	`

	var movies []repository.MovieSelectItem
	if err := r.db.SelectContext(ctx, &movies, query, userID, collectionID); err != nil {
		return nil, fmt.Errorf("list collection movies: %w", err)
	}

	return movies, nil
}

func (r *CollectionRepository) ListShows(ctx context.Context, userID int64, collectionID uuid.UUID) ([]repository.ShowSelectItem, error) {
	const query = `
		SELECT
			s.id,
			COALESCE(sm.name, s.name) as "name",
			s.year,
			s.fetch_source,
			s.path,
			COALESCE(sm.description, '') as "description",
			COALESCE(sm.medium_image_url, '') as "medium_image_url"
		FROM collection_shows cs
		JOIN collections c
			ON c.id = cs.collection_id
		JOIN shows s
			ON s.id = cs.show_id
		LEFT JOIN show_metadata sm
			ON sm.show_id = s.id AND s.fetch_source != 'multiple'
		WHERE c.user_id = $1 AND c.id = $2
		ORDER BY s.year DESC, LOWER(COALESCE(sm.name, s.name)) ASC
	`

	var shows []repository.ShowSelectItem
	if err := r.db.SelectContext(ctx, &shows, query, userID, collectionID); err != nil {
		return nil, fmt.Errorf("list collection shows: %w", err)
	}

	return shows, nil
}

func (r *CollectionRepository) ensureCollectionExistsForUser(ctx context.Context, collectionID uuid.UUID, userID int64) error {
	const query = `
		SELECT 1
		FROM collections
		WHERE id = $1 AND user_id = $2
		LIMIT 1
	`

	var exists int
	err := r.db.GetContext(ctx, &exists, query, collectionID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.ErrCollectionNotFound
		}
		return fmt.Errorf("ensure collection exists for user: %w", err)
	}

	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code == "23505"
}

var _ repository.CollectionRepository = (*CollectionRepository)(nil)
