package entity

import "github.com/google/uuid"

type EpisodeSubtitle struct {
	ID        uuid.UUID `db:"id"`
	EpisodeID uuid.UUID `db:"episode_id"`
	Path      string    `db:"path"`
	Name      string    `db:"name"`
	LangShort string    `db:"lang_short"`
	Lang      string    `db:"lang"`
}
