package entity

import "github.com/google/uuid"

type EpisodeMetadata struct {
	ID               uuid.UUID   `db:"id"`
	ShowID           uuid.UUID   `db:"show_id"`
	SeasonMetadataID uuid.UUID   `db:"season_metadata_id"`
	Url              string      `db:"url"`
	Name             string      `db:"name"`
	Number           int         `db:"number"`
	Summary          string      `db:"summary"`
	MediumImageUrl   string      `db:"medium_image_url"`
	OriginalImageUrl string      `db:"original_image_url"`
	FetchID          int         `db:"fetch_id"`
	FetchSource      FetchSource `db:"fetch_source"`
}
