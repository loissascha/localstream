package entity

import (
	"time"

	"github.com/google/uuid"
)

type CollectionMovie struct {
	ID           uuid.UUID `db:"id"`
	CollectionID uuid.UUID `db:"collection_id"`
	MovieID      uuid.UUID `db:"movie_id"`
	CreatedAt    time.Time `db:"created_at"`
}
