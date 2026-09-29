package entity

import (
	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/fetchsource"
)

type EpisodeMetadata struct {
	ID               uuid.UUID               `db:"id"`
	EpisodeID        uuid.UUID               `db:"episode_id"`
	Url              string                  `db:"url"`
	Name             string                  `db:"name"`
	Number           int                     `db:"number"`
	Summary          string                  `db:"summary"`
	MediumImageUrl   string                  `db:"medium_image_url"`
	OriginalImageUrl string                  `db:"original_image_url"`
	FetchID          int                     `db:"fetch_id"`
	FetchSource      fetchsource.FetchSource `db:"fetch_source"`
}
