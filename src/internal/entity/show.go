package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/fetchsource"
)

type Show struct {
	ID          uuid.UUID               `db:"id"`
	Name        string                  `db:"name"`
	Year        int                     `db:"year"`
	Path        string                  `db:"path"`
	CreatedAt   time.Time               `db:"created_at"`
	FetchSource fetchsource.FetchSource `db:"fetch_source"`
}
