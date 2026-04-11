package entity

import "github.com/google/uuid"

type SeasonMetadata struct {
	ID               uuid.UUID   `db:"id"`
	ShowID           uuid.UUID   `db:"show_id"`
	Url              string      `db:"url"`
	Number           int         `db:"number"`
	Summary          string      `db:"summary"`
	PremiereDate     string      `db:"premiere_date"`
	MediumImageUrl   string      `db:"medium_image_url"`
	OriginalImageUrl string      `db:"original_image_url"`
	FetchID          int         `db:"fetch_id"`
	FetchSource      FetchSource `db:"fetch_source"`
}
