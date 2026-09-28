package entity

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID          uuid.UUID   `db:"id"`
	Name        string      `db:"name"`
	Year        int         `db:"year"`
	Description string      `db:"description"`
	Path        string      `db:"path"`
	CreatedAt   time.Time   `db:"created_at"`
	FetchSource FetchSource `db:"fetch_source"`
}

type MovieVideoStream struct {
	ID            uuid.UUID         `db:"id"`
	MovieID       uuid.UUID         `db:"movie_id"`
	Index         int               `db:"stream_index"`
	Profile       string            `db:"profile"`
	Level         int               `db:"level"`
	CodecName     string            `db:"codec_name"`
	CodecLongName string            `db:"codec_long_name"`
	Duration      float64           `db:"duration"`
	Width         int               `db:"width"`
	Height        int               `db:"height"`
	Tags          map[string]string `db:"tags"`
	BitRate       int64             `db:"bit_rate"`
}

type MovieAudioStream struct {
	ID            uuid.UUID         `db:"id"`
	MovieID       uuid.UUID         `db:"movie_id"`
	Index         int               `db:"stream_index"`
	Profile       string            `db:"profile"`
	CodecName     string            `db:"codec_name"`
	CodecLongName string            `db:"codec_long_name"`
	Duration      float64           `db:"duration"`
	Tags          map[string]string `db:"tags"`
	BitRate       int64             `db:"bit_rate"`
}
