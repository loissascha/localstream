package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserWatchstate struct {
	ID        uuid.UUID `db:"id"`
	UserID    int64     `db:"user_id"`
	ShowID    uuid.UUID `db:"show_id"`
	SeasonID  uuid.UUID `db:"season_id"`
	EpisodeID uuid.UUID `db:"episode_id"`
	Position  float64   `db:"position"`
	Duration  float64   `db:"duration"`
	Finished  bool      `db:"finished"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
