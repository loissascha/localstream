package provider

import "github.com/google/uuid"

type SubtitleProviderResult struct {
	Name   string `json:"name"`
	Lang   string `json:"lang"`
	Author string `json:"author"`
	Url    string `json:"url"`
}

type SubtitleProvider interface {
	SearchMovie(name string) ([]SubtitleProviderResult, error)
	DownloadMovieSubtitle(movieId uuid.UUID, providerResult SubtitleProviderResult) error
}
