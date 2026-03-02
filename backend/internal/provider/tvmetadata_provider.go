package provider

import "github.com/loissascha/localstream/internal/parsers"

type ShowSearchResult struct {
	Score float64
	Show  ShowMetadata
}

type ShowMetadata struct {
	ID        int
	URL       string
	Name      string
	Genres    []string
	Premiered *string
	Image     *ShowImage
	Summary   *string
}

type ShowImage struct {
	Medium   string
	Original string
}

type TVMetadataProvider interface {
	SearchSeries(episodeInfo *parsers.EpisodeInfo) ([]ShowSearchResult, error)
}
