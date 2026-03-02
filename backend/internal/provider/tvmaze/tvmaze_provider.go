package tvmaze

import (
	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/parsers"
)

type TVMazeProvider struct {
}

func NewTVMazeProvider() *TVMazeProvider {
	return &TVMazeProvider{}
}

func (p *TVMazeProvider) SearchSeries(episodeInfo *parsers.EpisodeInfo) {
	if episodeInfo == nil {
		return
	}

	logger.Info(nil, "Search series {Name}", episodeInfo)
}
