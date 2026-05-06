package handler

import (
	"time"

	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type ShowInfo struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Year           int       `json:"year"`
	FetchSource    string    `json:"fetch_source"`
	Description    string    `json:"description"`
	MediumImageUrl string    `json:"medium_image_url"`
	CreatedAt      time.Time `json:"created_at"`
}

type ShowListResponse struct {
	Shows []ShowInfo `json:"shows"`
}

type ShowMetadataInfo struct {
	ID               string             `json:"id"`
	ShowID           string             `json:"show_id"`
	Name             string             `json:"name"`
	Url              string             `json:"url"`
	Description      string             `json:"description"`
	MediumImageUrl   string             `json:"medium_image_url"`
	OriginalImageUrl string             `json:"original_image_url"`
	FetchSource      entity.FetchSource `json:"fetch_source"`
}

func toShowMetadataInfo(m *entity.ShowMetadata) ShowMetadataInfo {
	return ShowMetadataInfo{
		ID:               encoders.EncodeUUID(m.ID),
		ShowID:           encoders.EncodeUUID(m.ShowID),
		Name:             m.Name,
		Url:              m.Url,
		Description:      m.Description,
		MediumImageUrl:   m.MediumImageUrl,
		OriginalImageUrl: m.OriginalImageUrl,
		FetchSource:      m.FetchSource,
	}
}

func toShowInfo(show *repository.ShowSelectItem) ShowInfo {
	return ShowInfo{
		ID:             encoders.EncodeUUID(show.ID),
		Name:           show.Name,
		Year:           show.Year,
		FetchSource:    string(show.FetchSource),
		Description:    show.Description,
		MediumImageUrl: show.MediumImageUrl,
		CreatedAt:      show.CreatedAt,
	}
}
