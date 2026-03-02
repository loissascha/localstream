package entity

import (
	"time"

	"github.com/google/uuid"
)

type Show struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}
