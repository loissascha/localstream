package provider

import (
	"context"

	"github.com/google/uuid"
)

type SubtitleProviderResult struct {
	Name   string `json:"name"`
	Lang   string `json:"lang"`
	Author string `json:"author"`
	Url    string `json:"url"`
}

type SubtitleProvider interface {
	SearchMovie(ctx context.Context, name string) ([]SubtitleProviderResult, error)
	DownloadMovieSubtitle(ctx context.Context, movieId uuid.UUID, providerResult SubtitleProviderResult) error
}
