package entity

import (
	"time"

	"github.com/google/uuid"
)

type MediaType string

const (
	MediaTypeEpisode MediaType = "episode"
	MediaTypeMovie   MediaType = "movie"
)

type UserUnifiedWatchstate struct {
	ID        uuid.UUID `db:"id"`
	UserID    int64     `db:"user_id"`
	MediaType MediaType `db:"media_type"`
	MediaID   uuid.UUID `db:"media_id"`
	Position  float64   `db:"position"`
	Duration  float64   `db:"duration"`
	Finished  bool      `db:"finished"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
