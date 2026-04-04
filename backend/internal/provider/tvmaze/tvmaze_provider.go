package tvmaze

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/provider"
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

func (p *TVMazeProvider) SearchShow(name string, year int) ([]provider.ShowSearchResult, error) {
	searchTerm := name
	if year != 0 {
		searchTerm = fmt.Sprintf("%s (%n)", name, year)
	}

	params := url.Values{}
	params.Add("q", searchTerm)
	encoded := params.Encode()

	fullUrl := "https://api.tvmaze.com/search/shows?" + encoded

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullUrl, nil)
	if err != nil {
		logger.Error(err, "Error with http request")
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err, "Error with http request 2")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err, "Error with http request 3")
		return nil, err
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("tvmaze search returned status %d", resp.StatusCode)
		logger.Error(err, "Non success status code")
		return nil, err
	}

	var searchResults []TVMazeSearchResult
	if err = json.Unmarshal(body, &searchResults); err != nil {
		logger.Error(err, "Error decoding tvmaze response")
		return nil, err
	}

	if len(searchResults) == 0 {
		logger.Info(nil, "No series found for {Name} ({Year})", name, year)
		return []provider.ShowSearchResult{}, nil
	}

	mappedResults := make([]provider.ShowSearchResult, 0, len(searchResults))
	for _, result := range searchResults {
		mappedResult := provider.ShowSearchResult{
			Score: result.Score,
			Show: provider.ShowMetadata{
				ID:        result.Show.ID,
				URL:       result.Show.URL,
				Name:      result.Show.Name,
				Genres:    result.Show.Genres,
				Premiered: result.Show.Premiered,
				Summary:   result.Show.Summary,
			},
		}

		if result.Show.Image != nil {
			mappedResult.Show.Image = &provider.ShowImage{
				Medium:   result.Show.Image.Medium,
				Original: result.Show.Image.Original,
			}
		}

		mappedResults = append(mappedResults, mappedResult)
	}

	firstResult := mappedResults[0]
	logger.Debug(nil, "Found {Count} results. Top match: {Name} ({ID}) with score {Score}", len(searchResults), firstResult.Show.Name, firstResult.Show.ID, firstResult.Score)

	return mappedResults, nil
}
