package handler

import (
	"time"

	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type MovieInfo struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Year             int       `json:"year"`
	Description      string    `json:"description"`
	FetchSource      string    `json:"fetch_source"`
	MediumImageUrl   string    `json:"medium_image_url"`
	BackdropImageUrl string    `json:"backdrop_image_url"`
	Position         float64   `json:"position"`
	Duration         float64   `json:"duration"`
	Finished         bool      `json:"finished"`
	Percentage       float64   `json:"percentage"`
	CreatedAt        time.Time `json:"created_at"`
}

type MovieMetadataInfo struct {
	ID               string             `json:"id"`
	MovieID          string             `json:"movie_id"`
	Name             string             `json:"name"`
	ReleaseYear      int                `json:"release_year"`
	Url              string             `json:"url"`
	Description      string             `json:"description"`
	MediumImageUrl   string             `json:"medium_image_url"`
	BackdropImageUrl string             `json:"backdrop_image_url"`
	FetchSource      entity.FetchSource `json:"fetch_source"`
}

type MovieSubtitleInfo struct {
	ID        string `json:"id"`
	MovieID   string `json:"movie_id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	LangShort string `json:"lang_short"`
	Lang      string `json:"lang"`
}

type MovieListResponse struct {
	Movies []MovieInfo `json:"movies"`
}

func toMovieSubtitleInfo(m *entity.MovieSubtitle) MovieSubtitleInfo {
	return MovieSubtitleInfo{
		ID:        encoders.EncodeUUID(m.ID),
		MovieID:   encoders.EncodeUUID(m.MovieID),
		Name:      m.Name,
		Path:      m.Path,
		LangShort: m.LangShort,
		Lang:      m.Lang,
	}
}

func toMovieInfo(m *repository.MovieSelectItem) MovieInfo {
	percent := 0.0
	if m.Duration > 0 {
		percent = (100 / m.Duration) * m.Position
	}
	if m.Finished {
		percent = 100
	}
	return MovieInfo{
		ID:               encoders.EncodeUUID(m.ID),
		Name:             m.Name,
		Year:             m.Year,
		Description:      m.Description,
		FetchSource:      string(m.FetchSource),
		MediumImageUrl:   m.MediumImageUrl,
		BackdropImageUrl: m.BackdropImageUrl,
		Duration:         m.Duration,
		Finished:         m.Finished,
		Position:         m.Position,
		Percentage:       percent,
		CreatedAt:        m.CreatedAt,
	}
}

func toMovieMetadataInfo(m *entity.MovieMetadata) MovieMetadataInfo {
	return MovieMetadataInfo{
		ID:               encoders.EncodeUUID(m.ID),
		MovieID:          encoders.EncodeUUID(m.MovieID),
		Name:             m.Name,
		ReleaseYear:      m.ReleaseYear,
		Url:              m.Url,
		Description:      m.Description,
		MediumImageUrl:   m.MediumImageUrl,
		BackdropImageUrl: m.BackdropImageUrl,
		FetchSource:      m.FetchSource,
	}
}
