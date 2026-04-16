package entity

import "github.com/google/uuid"

type MovieMetadata struct {
	ID               uuid.UUID   `db:"id"`
	MovieID          uuid.UUID   `db:"movie_id"`
	Name             string      `db:"name"`
	Url              string      `db:"url"`
	Description      string      `db:"description"`
	MediumImageUrl   string      `db:"medium_image_url"`
	BackdropImageUrl string      `db:"backdrop_image_url"`
	FetchSource      FetchSource `db:"fetch_source"`
}
