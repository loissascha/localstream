package handler

import (
	"time"

	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
)

type AnyInfoStruct interface{}

type SearchResponse struct {
	Shows  []ShowInfo  `json:"shows"`
	Movies []MovieInfo `json:"movies"`
}

type WatchstateMovieResponse struct {
	ID         string    `json:"id"`
	MovieID    string    `json:"movie_id"`
	Position   float64   `json:"position"`
	Duration   float64   `json:"duration"`
	Finished   bool      `json:"finished"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Percentage float64   `json:"percentage"`
	MovieInfo  MovieInfo `json:"movie_info"`
}

type WatchstateResponse struct {
	ID          string      `json:"id"`
	ShowID      string      `json:"show_id"`
	ShowInfo    ShowInfo    `json:"show_info"`
	SeasonID    string      `json:"season_id"`
	SeasonInfo  SeasonInfo  `json:"season_info"`
	EpisodeID   string      `json:"episode_id"`
	EpisodeInfo EpisodeInfo `json:"episode_info"`
	Position    float64     `json:"position"`
	Duration    float64     `json:"duration"`
	Finished    bool        `json:"finished"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Percentage  float64     `json:"percentage"`
}

type WatchstateInfo struct {
	Position   float64 `json:"position"`
	Duration   float64 `json:"duration"`
	Percentage float64 `json:"percentage"`
	Finished   bool    `json:"finished"`
}

type SeasonInfo struct {
	ID     string `json:"id"`
	Number int    `json:"number"`
}

type SeasonListResponse struct {
	Seasons []SeasonInfo `json:"seasons"`
}

type LibraryListItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	LibraryType string `json:"library_type"`
}

type LibraryListResponse struct {
	Libraries []LibraryListItem `json:"libraries"`
}

type SeasonMetadataInfo struct {
	ID               string             `json:"id"`
	SeasonID         string             `json:"season_id"`
	Url              string             `json:"url"`
	Number           int                `json:"number"`
	Summary          string             `json:"summary"`
	PremiereDate     string             `json:"premiere_date"`
	MediumImageUrl   string             `json:"medium_image_url"`
	OriginalImageUrl string             `json:"original_image_url"`
	FetchID          int                `json:"fetch_id"`
	FetchSource      entity.FetchSource `json:"fetch_source"`
}

func toSeasonMetadataInfo(m *entity.SeasonMetadata) SeasonMetadataInfo {
	return SeasonMetadataInfo{
		ID:               encoders.EncodeUUID(m.ID),
		SeasonID:         encoders.EncodeUUID(m.SeasonID),
		Url:              m.Url,
		Number:           m.Number,
		Summary:          m.Summary,
		PremiereDate:     m.PremiereDate,
		MediumImageUrl:   m.MediumImageUrl,
		OriginalImageUrl: m.OriginalImageUrl,
		FetchID:          m.FetchID,
		FetchSource:      m.FetchSource,
	}
}

func toLibraryListItem(l *entity.Library) LibraryListItem {
	return LibraryListItem{
		ID:          encoders.EncodeUUID(l.ID),
		Name:        l.Name,
		Path:        l.Path,
		LibraryType: string(l.LibraryType),
	}
}

func toWatchstateInfoMovie(watchstate *entity.UserMovieWatchstate) WatchstateInfo {
	percent := 0.0
	if watchstate.Duration > 0 {
		percent = (100 / watchstate.Duration) * watchstate.Position
	}
	if watchstate.Finished {
		percent = 100
	}
	return WatchstateInfo{
		Position:   watchstate.Position,
		Duration:   watchstate.Duration,
		Percentage: percent,
		Finished:   watchstate.Finished,
	}
}

func toWatchstateInfo(watchstate *entity.UserWatchstate) WatchstateInfo {
	percent := 0.0
	if watchstate.Duration > 0 {
		percent = (100 / watchstate.Duration) * watchstate.Position
	}
	if watchstate.Finished {
		percent = 100
	}
	return WatchstateInfo{
		Position:   watchstate.Position,
		Duration:   watchstate.Duration,
		Percentage: percent,
		Finished:   watchstate.Finished,
	}
}

func toWatchstateMovieResponse(watchstate entity.UserMovieWatchstate, movieInfo MovieInfo) WatchstateMovieResponse {
	percent := 0.0
	if watchstate.Duration > 0 {
		percent = (100 / watchstate.Duration) * watchstate.Position
	}
	if watchstate.Finished {
		percent = 100
	}
	return WatchstateMovieResponse{
		ID:         encoders.EncodeUUID(watchstate.ID),
		MovieID:    encoders.EncodeUUID(watchstate.MovieID),
		Position:   watchstate.Position,
		Duration:   watchstate.Duration,
		Finished:   watchstate.Finished,
		Percentage: percent,
		MovieInfo:  movieInfo,
	}
}

func toWatchstateResponse(watchstate entity.UserWatchstate, infos ...AnyInfoStruct) WatchstateResponse {
	var showInfo ShowInfo
	var seasonInfo SeasonInfo
	var episodeInfo EpisodeInfo
	for _, info := range infos {
		eI, ok := info.(EpisodeInfo)
		if ok {
			episodeInfo = eI
			continue
		}
		shI, ok := info.(ShowInfo)
		if ok {
			showInfo = shI
			continue
		}
		seI, ok := info.(SeasonInfo)
		if ok {
			seasonInfo = seI
			continue
		}
	}
	percent := 0.0
	if watchstate.Duration > 0 {
		percent = (100 / watchstate.Duration) * watchstate.Position
	}
	if watchstate.Finished {
		percent = 100
	}
	return WatchstateResponse{
		ID:          encoders.EncodeUUID(watchstate.ID),
		ShowID:      encoders.EncodeUUID(watchstate.ShowID),
		ShowInfo:    showInfo,
		SeasonID:    encoders.EncodeUUID(watchstate.SeasonID),
		SeasonInfo:  seasonInfo,
		EpisodeID:   encoders.EncodeUUID(watchstate.EpisodeID),
		EpisodeInfo: episodeInfo,
		Percentage:  percent,
		Position:    watchstate.Position,
		Duration:    watchstate.Duration,
		Finished:    watchstate.Finished,
		CreatedAt:   watchstate.CreatedAt,
		UpdatedAt:   watchstate.UpdatedAt,
	}
}
