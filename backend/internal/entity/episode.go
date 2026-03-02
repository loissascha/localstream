package entity

import (
	"time"

	"github.com/google/uuid"
)

type Episode struct {
	ID          uuid.UUID   `db:"id"`
	SeasonID    uuid.UUID   `db:"season_id"`
	Name        string      `db:"name"`
	Path        string      `db:"path"`
	CreatedAt   time.Time   `db:"created_at"`
	FetchSource FetchSource `db:"fetch_source"`
}
