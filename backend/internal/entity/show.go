package entity

import (
	"time"

	"github.com/google/uuid"
)

type Show struct {
	ID          uuid.UUID   `db:"id"`
	Name        string      `db:"name"`
	Path        string      `db:"path"`
	CreatedAt   time.Time   `db:"created_at"`
	FetchSource FetchSource `db:"fetch_source"`
}

type FetchSource string

const (
	FetchSourceNone   FetchSource = "none"
	FetchSourceTMDB   FetchSource = "tmdb"
	FetchSourceTVMaze FetchSource = "tvmaze"
)

func (t FetchSource) IsValid() bool {
	switch t {
	case FetchSourceTMDB, FetchSourceTVMaze:
		return true
	default:
		return false
	}
}
