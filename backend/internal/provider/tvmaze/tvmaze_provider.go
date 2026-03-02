package tvmaze

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/parsers"
)

type TVMazeProvider struct {
}

type TVMazeSearchResult struct {
	Score float64    `json:"score"`
	Show  TVMazeShow `json:"show"`
}

type TVMazeShow struct {
	ID        int              `json:"id"`
	URL       string           `json:"url"`
	Name      string           `json:"name"`
	Genres    []string         `json:"genres"`
	Premiered *string          `json:"premiered"`
	Image     *TVMazeShowImage `json:"image"`
	Summary   *string          `json:"summary"`
}

type TVMazeShowImage struct {
	Medium   string `json:"medium"`
	Original string `json:"original"`
}

func NewTVMazeProvider() *TVMazeProvider {
	return &TVMazeProvider{}
}

func (p *TVMazeProvider) SearchSeries(episodeInfo *parsers.EpisodeInfo) {
	if episodeInfo == nil {
		return
	}

	logger.Info(nil, "Search series {Name}", episodeInfo)

	params := url.Values{}
	params.Add("q", episodeInfo.Series)
	encoded := params.Encode()

	fullUrl := "https://api.tvmaze.com/search/shows?" + encoded

	logger.Info(nil, "URL: {Url}", fullUrl)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullUrl, nil)
	if err != nil {
		logger.Error(err, "Error with http request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err, "Error with http request 2")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err, "Error with http request 3")
		return
	}

	if resp.StatusCode != 200 {
		logger.Error(nil, "Non success status code")
		return
	}

	var searchResults []TVMazeSearchResult
	if err = json.Unmarshal(body, &searchResults); err != nil {
		logger.Error(err, "Error decoding tvmaze response")
		return
	}

	if len(searchResults) == 0 {
		logger.Info(nil, "No series found for {Name}", episodeInfo.Series)
		return
	}

	firstResult := searchResults[0]
	logger.Debug(nil, "Found {Count} results. Top match: {Name} ({ID}) with score {Score}", len(searchResults), firstResult.Show.Name, firstResult.Show.ID, firstResult.Score)
}
