package tmdb

import (
	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/parsers"
)

type TMDBProvider struct {
}

func NewTMDBProvider() *TMDBProvider {
	return &TMDBProvider{}
}

func (p *TMDBProvider) SearchSeries(episodeInfo *parsers.EpisodeInfo) {
	if episodeInfo == nil {
		return
	}

	logger.Info(nil, "Search series {Name}", episodeInfo)
}
