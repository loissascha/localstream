package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserMovieWatchstate struct {
	ID        uuid.UUID `db:"id"`
	UserID    int64     `db:"user_id"`
	MovieID   uuid.UUID `db:"movie_id"`
	Position  float64   `db:"position"`
	Duration  float64   `db:"duration"`
	Finished  bool      `db:"finished"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
