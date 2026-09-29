package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/fetchsource"
)

type Season struct {
	ID          uuid.UUID               `db:"id"`
	ShowID      uuid.UUID               `db:"show_id"`
	Number      int                     `db:"number"`
	Path        string                  `db:"path"`
	CreatedAt   time.Time               `db:"created_at"`
	FetchSource fetchsource.FetchSource `db:"fetch_source"`
}
