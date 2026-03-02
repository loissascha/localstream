package entity

import (
	"time"

	"github.com/google/uuid"
)

type Library struct {
	ID          uuid.UUID   `db:"id"`
	Name        string      `db:"name"`
	Path        string      `db:"path"`
	LibraryType LibraryType `db:"library_type"`
	CreatedAt   time.Time   `db:"created_at"`
}

type LibraryType string

const (
	LibraryTypeMovies LibraryType = "movies"
	LibraryTypeShows  LibraryType = "shows"
)

func (t LibraryType) IsValid() bool {
	switch t {
	case LibraryTypeMovies, LibraryTypeShows:
		return true
	default:
		return false
	}
}
