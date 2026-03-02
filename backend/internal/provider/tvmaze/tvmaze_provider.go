package tvmaze

import (
	"context"
	"net/http"
	"strings"
	"time"

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

	var urlBuilder strings.Builder
	urlBuilder.WriteString("https://api.tvmaze.com/search/shows?q=")
	urlBuilder.WriteString(episodeInfo.Series)
	url := urlBuilder.String()

	logger.Info(nil, "URL: {Url}", url)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	http.NewRequestWithContext()
}
