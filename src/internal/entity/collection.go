package entity

import (
	"time"

	"github.com/google/uuid"
)

type Collection struct {
	ID        uuid.UUID `db:"id"`
	UserID    int64     `db:"user_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
