package entity

import (
	"time"

	"github.com/google/uuid"
)

type Library struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Path      string    `db:"path"`
	CreatedAt time.Time `db:"created_at"`
}
