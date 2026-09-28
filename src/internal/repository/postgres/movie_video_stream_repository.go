package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type MovieVideoStreamRepository struct {
	db *sqlx.DB
}

func NewMovieVideoStreamRepository(db *sqlx.DB) *MovieVideoStreamRepository {
	return &MovieVideoStreamRepository{db: db}
}

func (r *MovieVideoStreamRepository) Create(ctx context.Context, stream *entity.MovieVideoStream) error {
	tags, err := marshalStreamTags(stream.Tags)
	if err != nil {
		return fmt.Errorf("encode movie video stream tags: %w", err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("generate movie video stream id: %w", err)
	}

	const query = `
		INSERT INTO movie_video_streams (
			id, movie_id, stream_index, profile, level, codec_name, codec_long_name,
			duration, width, height, tags, bit_rate
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err = r.db.ExecContext(ctx, query,
		id, stream.MovieID, stream.Index, stream.Profile, stream.Level,
		stream.CodecName, stream.CodecLongName, stream.Duration, stream.Width,
		stream.Height, string(tags), stream.BitRate,
	)
	if err != nil {
		return fmt.Errorf("create movie video stream: %w", err)
	}

	stream.ID = id
	return nil
}

func (r *MovieVideoStreamRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.MovieVideoStream, error) {
	const query = `
		SELECT id, movie_id, stream_index, profile, level, codec_name, codec_long_name,
			duration, width, height, tags, bit_rate
		FROM movie_video_streams
		WHERE id = $1
		LIMIT 1
	`

	stream, err := scanMovieVideoStream(r.db.QueryRowxContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get movie video stream by id: %w", err)
	}
	return stream, nil
}

func (r *MovieVideoStreamRepository) ListByMovieID(ctx context.Context, movieID uuid.UUID) ([]entity.MovieVideoStream, error) {
	const query = `
		SELECT id, movie_id, stream_index, profile, level, codec_name, codec_long_name,
			duration, width, height, tags, bit_rate
		FROM movie_video_streams
		WHERE movie_id = $1
		ORDER BY stream_index ASC
	`

	rows, err := r.db.QueryxContext(ctx, query, movieID)
	if err != nil {
		return nil, fmt.Errorf("list movie video streams by movie id: %w", err)
	}
	defer rows.Close()

	streams := make([]entity.MovieVideoStream, 0)
	for rows.Next() {
		stream, err := scanMovieVideoStream(rows)
		if err != nil {
			return nil, fmt.Errorf("scan movie video stream: %w", err)
		}
		streams = append(streams, *stream)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list movie video streams by movie id: %w", err)
	}
	return streams, nil
}

func (r *MovieVideoStreamRepository) Update(ctx context.Context, stream *entity.MovieVideoStream) error {
	tags, err := marshalStreamTags(stream.Tags)
	if err != nil {
		return fmt.Errorf("encode movie video stream tags: %w", err)
	}

	const query = `
		UPDATE movie_video_streams
		SET movie_id = $1, stream_index = $2, profile = $3, level = $4,
			codec_name = $5, codec_long_name = $6, duration = $7, width = $8,
			height = $9, tags = $10, bit_rate = $11
		WHERE id = $12
	`
	result, err := r.db.ExecContext(ctx, query,
		stream.MovieID, stream.Index, stream.Profile, stream.Level,
		stream.CodecName, stream.CodecLongName, stream.Duration, stream.Width,
		stream.Height, string(tags), stream.BitRate, stream.ID,
	)
	if err != nil {
		return fmt.Errorf("update movie video stream: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update movie video stream rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return repository.ErrMovieVideoStreamNotFound
	}
	return nil
}

func (r *MovieVideoStreamRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM movie_video_streams WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete movie video stream: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete movie video stream rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return repository.ErrMovieVideoStreamNotFound
	}
	return nil
}

type movieVideoStreamRow interface {
	Scan(dest ...any) error
}

func scanMovieVideoStream(row movieVideoStreamRow) (*entity.MovieVideoStream, error) {
	var stream entity.MovieVideoStream
	var tags []byte
	if err := row.Scan(
		&stream.ID, &stream.MovieID, &stream.Index, &stream.Profile, &stream.Level,
		&stream.CodecName, &stream.CodecLongName, &stream.Duration, &stream.Width,
		&stream.Height, &tags, &stream.BitRate,
	); err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		stream.Tags = make(map[string]string)
	} else if err := json.Unmarshal(tags, &stream.Tags); err != nil {
		return nil, fmt.Errorf("decode tags: %w", err)
	}
	return &stream, nil
}

func marshalStreamTags(tags map[string]string) ([]byte, error) {
	if tags == nil {
		tags = map[string]string{}
	}
	return json.Marshal(tags)
}

var _ repository.MovieVideoStreamRepository = (*MovieVideoStreamRepository)(nil)
