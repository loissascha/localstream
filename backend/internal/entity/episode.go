package entity

import (
	"time"

	"github.com/google/uuid"
)

type Episode struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Path      string    `db:"path"`
	CreatedAt time.Time `db:"created_at"`
}
