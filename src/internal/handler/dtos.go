package handler

import (
	"time"

	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type AnyInfoStruct interface{}

type MovieInfo struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Year             int    `json:"year"`
	Description      string `json:"description"`
	FetchSource      string `json:"fetch_source"`
	MediumImageUrl   string `db:"medium_image_url"`
	BackdropImageUrl string `db:"backdrop_image_url"`
}

type MovieListResponse struct {
	Movies []MovieInfo `json:"movies"`
}

type ShowInfo struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Year           int    `json:"year"`
	FetchSource    string `json:"fetch_source"`
	Description    string `json:"description"`
	MediumImageUrl string `json:"medium_image_url"`
}

type ShowListResponse struct {
	Shows []ShowInfo `json:"shows"`
}

type SearchResponse struct {
	Shows  []ShowInfo  `json:"shows"`
	Movies []MovieInfo `json:"movies"`
}

type EpisodeInfo struct {
	ID         string         `json:"id"`
	SeasonID   string         `json:"season_id"`
	Number     int            `json:"number"`
	Watchstate WatchstateInfo `json:"watchstate"`
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

type EpisodeListResponse struct {
	Episodes []EpisodeInfo `json:"episodes"`
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

type MovieMetadataInfo struct {
	ID               string             `json:"id"`
	MovieID          string             `json:"movie_id"`
	Name             string             `json:"name"`
	Url              string             `json:"url"`
	Description      string             `json:"description"`
	MediumImageUrl   string             `json:"medium_image_url"`
	BackdropImageUrl string             `json:"backdrop_image_url"`
	FetchSource      entity.FetchSource `json:"fetch_source"`
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

type EpisodeMetadataInfo struct {
	ID               string             `json:"id"`
	EpisodeID        string             `json:"episode_id"`
	Url              string             `json:"url"`
	Name             string             `json:"name"`
	Number           int                `json:"number"`
	Summary          string             `json:"summary"`
	MediumImageUrl   string             `json:"medium_image_url"`
	OriginalImageUrl string             `json:"original_image_url"`
	FetchID          int                `json:"fetch_id"`
	FetchSource      entity.FetchSource `json:"fetch_source"`
}

func toMovieMetadataInfo(m *entity.MovieMetadata) MovieMetadataInfo {
	return MovieMetadataInfo{
		ID:               encoders.EncodeUUID(m.ID),
		MovieID:          encoders.EncodeUUID(m.MovieID),
		Name:             m.Name,
		Url:              m.Url,
		Description:      m.Description,
		MediumImageUrl:   m.MediumImageUrl,
		BackdropImageUrl: m.BackdropImageUrl,
		FetchSource:      m.FetchSource,
	}
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

func toEpisodeMetadataInfo(m *entity.EpisodeMetadata) EpisodeMetadataInfo {
	return EpisodeMetadataInfo{
		ID:               encoders.EncodeUUID(m.ID),
		EpisodeID:        encoders.EncodeUUID(m.EpisodeID),
		Url:              m.Url,
		Name:             m.Name,
		Number:           m.Number,
		Summary:          m.Summary,
		MediumImageUrl:   m.MediumImageUrl,
		OriginalImageUrl: m.OriginalImageUrl,
		FetchID:          m.FetchID,
		FetchSource:      m.FetchSource,
	}
}

func toMovieInfo(m *repository.MovieSelectItem) MovieInfo {
	return MovieInfo{
		ID:               encoders.EncodeUUID(m.ID),
		Name:             m.Name,
		Year:             m.Year,
		Description:      m.Description,
		FetchSource:      string(m.FetchSource),
		MediumImageUrl:   m.MediumImageUrl,
		BackdropImageUrl: m.BackdropImageUrl,
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

func toEpisodeInfo(episode *entity.Episode, infos ...AnyInfoStruct) EpisodeInfo {
	var watchstateInfo WatchstateInfo
	for _, info := range infos {
		i, ok := info.(WatchstateInfo)
		if ok {
			watchstateInfo = i
		}
	}

	return EpisodeInfo{
		ID:         encoders.EncodeUUID(episode.ID),
		SeasonID:   encoders.EncodeUUID(episode.SeasonID),
		Number:     episode.Number,
		Watchstate: watchstateInfo,
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
