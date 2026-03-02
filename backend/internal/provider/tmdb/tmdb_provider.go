package tmdb

import (
	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/parsers"
	"github.com/loissascha/localstream/internal/provider"
)

type TMDBProvider struct {
}

func NewTMDBProvider() *TMDBProvider {
	return &TMDBProvider{}
}

func (p *TMDBProvider) SearchSeries(episodeInfo *parsers.EpisodeInfo) ([]provider.ShowSearchResult, error) {
	if episodeInfo == nil {
		return nil, nil
	}

	logger.Info(nil, "Search series {Name}", episodeInfo)

	return []provider.ShowSearchResult{}, nil
}
