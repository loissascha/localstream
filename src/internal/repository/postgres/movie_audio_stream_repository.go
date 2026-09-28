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

type MovieAudioStreamRepository struct {
	db *sqlx.DB
}

func NewMovieAudioStreamRepository(db *sqlx.DB) *MovieAudioStreamRepository {
	return &MovieAudioStreamRepository{db: db}
}

func (r *MovieAudioStreamRepository) Create(ctx context.Context, stream *entity.MovieAudioStream) error {
	tags, err := marshalStreamTags(stream.Tags)
	if err != nil {
		return fmt.Errorf("encode movie audio stream tags: %w", err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("generate movie audio stream id: %w", err)
	}

	const query = `
		INSERT INTO movie_audio_streams (
			id, movie_id, stream_index, profile, codec_name, codec_long_name,
			duration, tags, bit_rate
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = r.db.ExecContext(ctx, query,
		id, stream.MovieID, stream.Index, stream.Profile, stream.CodecName,
		stream.CodecLongName, stream.Duration, string(tags), stream.BitRate,
	)
	if err != nil {
		return fmt.Errorf("create movie audio stream: %w", err)
	}

	stream.ID = id
	return nil
}

func (r *MovieAudioStreamRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.MovieAudioStream, error) {
	const query = `
		SELECT id, movie_id, stream_index, profile, codec_name, codec_long_name,
			duration, tags, bit_rate
		FROM movie_audio_streams
		WHERE id = $1
		LIMIT 1
	`

	stream, err := scanMovieAudioStream(r.db.QueryRowxContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get movie audio stream by id: %w", err)
	}
	return stream, nil
}

func (r *MovieAudioStreamRepository) ListByMovieID(ctx context.Context, movieID uuid.UUID) ([]entity.MovieAudioStream, error) {
	const query = `
		SELECT id, movie_id, stream_index, profile, codec_name, codec_long_name,
			duration, tags, bit_rate
		FROM movie_audio_streams
		WHERE movie_id = $1
		ORDER BY stream_index ASC
	`

	rows, err := r.db.QueryxContext(ctx, query, movieID)
	if err != nil {
		return nil, fmt.Errorf("list movie audio streams by movie id: %w", err)
	}
	defer rows.Close()

	streams := make([]entity.MovieAudioStream, 0)
	for rows.Next() {
		stream, err := scanMovieAudioStream(rows)
		if err != nil {
			return nil, fmt.Errorf("scan movie audio stream: %w", err)
		}
		streams = append(streams, *stream)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list movie audio streams by movie id: %w", err)
	}
	return streams, nil
}

func (r *MovieAudioStreamRepository) Update(ctx context.Context, stream *entity.MovieAudioStream) error {
	tags, err := marshalStreamTags(stream.Tags)
	if err != nil {
		return fmt.Errorf("encode movie audio stream tags: %w", err)
	}

	const query = `
		UPDATE movie_audio_streams
		SET movie_id = $1, stream_index = $2, profile = $3, codec_name = $4,
			codec_long_name = $5, duration = $6, tags = $7, bit_rate = $8
		WHERE id = $9
	`
	result, err := r.db.ExecContext(ctx, query,
		stream.MovieID, stream.Index, stream.Profile, stream.CodecName,
		stream.CodecLongName, stream.Duration, string(tags), stream.BitRate,
		stream.ID,
	)
	if err != nil {
		return fmt.Errorf("update movie audio stream: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update movie audio stream rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return repository.ErrMovieAudioStreamNotFound
	}
	return nil
}

func (r *MovieAudioStreamRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM movie_audio_streams WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete movie audio stream: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete movie audio stream rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return repository.ErrMovieAudioStreamNotFound
	}
	return nil
}

func scanMovieAudioStream(row movieVideoStreamRow) (*entity.MovieAudioStream, error) {
	var stream entity.MovieAudioStream
	var tags []byte
	if err := row.Scan(
		&stream.ID, &stream.MovieID, &stream.Index, &stream.Profile,
		&stream.CodecName, &stream.CodecLongName, &stream.Duration, &tags,
		&stream.BitRate,
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

var _ repository.MovieAudioStreamRepository = (*MovieAudioStreamRepository)(nil)
