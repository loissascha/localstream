package handler

import (
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
	ID     string `json:"id"`
	Number int    `json:"number"`
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

func toEpisodeInfo(episode *entity.Episode) EpisodeInfo {
	return EpisodeInfo{
		ID:     encoders.EncodeUUID(episode.ID),
		Number: episode.Number,
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
