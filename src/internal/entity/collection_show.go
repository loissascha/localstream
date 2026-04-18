package entity

import (
	"time"

	"github.com/google/uuid"
)

type CollectionShow struct {
	ID           uuid.UUID `db:"id"`
	CollectionID uuid.UUID `db:"collection_id"`
	ShowID       uuid.UUID `db:"show_id"`
	CreatedAt    time.Time `db:"created_at"`
}
