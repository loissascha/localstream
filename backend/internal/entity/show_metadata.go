package entity

import "github.com/google/uuid"

type ShowMetadata struct {
	ID               uuid.UUID   `db:"id"`
	ShowID           uuid.UUID   `db:"show_id"`
	Name             string      `db:"name"`
	Url              string      `db:"url"`
	Description      string      `db:"description"`
	MediumImageUrl   string      `db:"medium_image_url"`
	OriginalImageUrl string      `db:"original_image_url"`
	FetchSource      FetchSource `db:"fetch_source"`
}
