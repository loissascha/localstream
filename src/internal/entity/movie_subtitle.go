package entity

import "github.com/google/uuid"

type MovieSubtitle struct {
	ID        uuid.UUID `db:"id"`
	MovieID   uuid.UUID `db:"movie_id"`
	Path      string    `db:"path"`
	Name      string    `db:"name"`
	LangShort string    `db:"lang_short"`
	Lang      string    `db:"lang"`
}
