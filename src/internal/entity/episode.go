package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/fetchsource"
)

type Episode struct {
	ID          uuid.UUID               `db:"id"`
	SeasonID    uuid.UUID               `db:"season_id"`
	Number      int                     `db:"number"`
	Path        string                  `db:"path"`
	CreatedAt   time.Time               `db:"created_at"`
	FetchSource fetchsource.FetchSource `db:"fetch_source"`
}
