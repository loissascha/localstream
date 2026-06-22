package handler

import (
	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type EpisodeInfo struct {
	ID               string         `json:"id"`
	SeasonID         string         `json:"season_id"`
	Number           int            `json:"number"`
	Watchstate       WatchstateInfo `json:"watchstate"`
	Name             string         `json:"name"`
	Summary          string         `json:"summary"`
	MediumImageUrl   string         `json:"medium_image_url"`
	OriginalImageUrl string         `json:"original_image_url"`
	FetchID          int            `json:"fetch_id"`
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

type EpisodeListResponse struct {
	Episodes []EpisodeInfo `json:"episodes"`
}

func toSubtitleInfoEpisode(m *entity.EpisodeSubtitle) SubtitleInfo {
	return SubtitleInfo{
		ID:        encoders.EncodeUUID(m.ID),
		Name:      m.Name,
		Path:      m.Path,
		LangShort: m.LangShort,
		Lang:      m.Lang,
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

func toEpisodeInfo(episode *repository.EpisodeWithMetadata, infos ...AnyInfoStruct) EpisodeInfo {
	var watchstateInfo WatchstateInfo
	for _, info := range infos {
		i, ok := info.(WatchstateInfo)
		if ok {
			watchstateInfo = i
		}
	}

	return EpisodeInfo{
		ID:               encoders.EncodeUUID(episode.ID),
		SeasonID:         encoders.EncodeUUID(episode.SeasonID),
		Number:           episode.Number,
		Watchstate:       watchstateInfo,
		Name:             episode.Name,
		Summary:          episode.Summary,
		MediumImageUrl:   episode.MediumImageUrl,
		OriginalImageUrl: episode.OriginalImageUrl,
		FetchID:          episode.FetchID,
	}
}
