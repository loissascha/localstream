package handler

import (
	"time"

	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
)

type CollectionInfo struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Movies    []MovieInfo `json:"movies"`
	Shows     []ShowInfo  `json:"shows"`
}

type CollectionInfoOption func(*CollectionInfo)

type CollectionListResponse struct {
	Collections []CollectionInfo `json:"collections"`
}

type CollectionDetailResponse struct {
	Collection CollectionInfo `json:"collection"`
	Movies     []MovieInfo    `json:"movies"`
	Shows      []ShowInfo     `json:"shows"`
}

func toCollectionInfo(c *entity.Collection, options ...CollectionInfoOption) CollectionInfo {
	ci := CollectionInfo{
		ID:        encoders.EncodeUUID(c.ID),
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Movies:    []MovieInfo{},
		Shows:     []ShowInfo{},
	}

	for _, opt := range options {
		opt(&ci)
	}

	return ci
}

func withCollectionMovies(movies []MovieInfo) CollectionInfoOption {
	return func(ci *CollectionInfo) {
		ci.Movies = movies
	}
}

func withCollectionShows(shows []ShowInfo) CollectionInfoOption {
	return func(ci *CollectionInfo) {
		ci.Shows = shows
	}
}
