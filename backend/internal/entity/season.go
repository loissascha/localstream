package entity

import (
	"time"

	"github.com/google/uuid"
)

type Season struct {
	ID          uuid.UUID   `db:"id"`
	Name        string      `db:"name"`
	CreatedAt   time.Time   `db:"created_at"`
	FetchSource FetchSource `db:"fetch_source"`
}
