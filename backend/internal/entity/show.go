package entity

import (
	"time"

	"github.com/google/uuid"
)

type Show struct {
	ID          uuid.UUID   `db:"id"`
	Name        string      `db:"name"`
	CreatedAt   time.Time   `db:"created_at"`
	FetchSource FetchSource `db:"fetch_source"`
}

type FetchSource string

const (
	FetchSourceTMDB   FetchSource = "tmdb"
	FetchSourceTVMaze FetchSource = "tvmaze"
	FetchSourceNone   FetchSource = "none"
)

func (t FetchSource) IsValid() bool {
	switch t {
	case FetchSourceTMDB, FetchSourceTVMaze:
		return true
	default:
		return false
	}
}
