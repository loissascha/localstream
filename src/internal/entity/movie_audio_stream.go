package entity

import "github.com/google/uuid"

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
