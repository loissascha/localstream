package entity

import (
	"time"

	"github.com/google/uuid"
)

type Season struct {
	ID          uuid.UUID   `db:"id"`
	ShowID      uuid.UUID   `db:"show_id"`
	Name        string      `db:"name"`
	Path        string      `db:"path"`
	CreatedAt   time.Time   `db:"created_at"`
	FetchSource FetchSource `db:"fetch_source"`
}
