package entity

import (
	"time"

	"github.com/google/uuid"
)

type Show struct {
	ID          uuid.UUID   `db:"id"`
	Name        string      `db:"name"`
	Year        int         `db:"year"`
	Description string      `db:"description"`
	Path        string      `db:"path"`
	CreatedAt   time.Time   `db:"created_at"`
	FetchSource FetchSource `db:"fetch_source"`
}
