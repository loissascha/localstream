package entity

import "github.com/google/uuid"

type MovieVideoStream struct {
	ID            uuid.UUID         `db:"id"`
	MovieID       uuid.UUID         `db:"movie_id"`
	Container     string            `db:"container"`
	Index         int               `db:"stream_index"`
	Profile       string            `db:"profile"`
	Level         int               `db:"level"`
	CodecName     string            `db:"codec_name"`
	CodecLongName string            `db:"codec_long_name"`
	PixelFormat   string            `db:"pixel_format"`
	Duration      float64           `db:"duration"`
	Width         int               `db:"width"`
	Height        int               `db:"height"`
	Tags          map[string]string `db:"tags"`
	BitRate       int64             `db:"bit_rate"`
}
