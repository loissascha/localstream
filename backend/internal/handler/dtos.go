package handler

import (
	"time"

	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
)

type AnyInfoStruct interface{}

type ShowInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Year        int    `json:"year"`
	Description string `json:"description"`
}

type ShowListResponse struct {
	Shows []ShowInfo `json:"shows"`
}

type EpisodeInfo struct {
	ID         string         `json:"id"`
	Number     int            `json:"number"`
	Watchstate WatchstateInfo `json:"watchstate"`
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
		Number:     episode.Number,
		Watchstate: watchstateInfo,
	}
}

func toShowInfo(show *entity.Show) ShowInfo {
	return ShowInfo{
		ID:          encoders.EncodeUUID(show.ID),
		Name:        show.Name,
		Year:        show.Year,
		Description: show.Description,
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
	return WatchstateResponse{
		ID:          encoders.EncodeUUID(watchstate.ID),
		ShowID:      encoders.EncodeUUID(watchstate.ShowID),
		ShowInfo:    showInfo,
		SeasonID:    encoders.EncodeUUID(watchstate.SeasonID),
		SeasonInfo:  seasonInfo,
		EpisodeID:   encoders.EncodeUUID(watchstate.EpisodeID),
		EpisodeInfo: episodeInfo,
		Position:    watchstate.Position,
		Duration:    watchstate.Duration,
		Finished:    watchstate.Finished,
		CreatedAt:   watchstate.CreatedAt,
		UpdatedAt:   watchstate.UpdatedAt,
	}
}
