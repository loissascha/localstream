package tvmaze

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	Ended     *string          `json:"ended"`
	Image     *TVMazeShowImage `json:"image"`
	Summary   *string          `json:"summary"`
	Language  *string          `json:"language"`
}

type TVMazeSeason struct {
}

type TVMazeShowImage struct {
	Medium   string `json:"medium"`
	Original string `json:"original"`
}

func NewTVMazeProvider() *TVMazeProvider {
	return &TVMazeProvider{}
}

func (p *TVMazeProvider) SearchSeasons(showId int) ([]provider.SeasonMetadata, error) {
	fullUrl := fmt.Sprintf("https://api.tvmaze.com/shows/%d/seasons", showId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("tvmaze search returned status %d", resp.StatusCode)
		return nil, err
	}

	var searchResults []TVMazeSeason
	if err = json.Unmarshal(body, &searchResults); err != nil {
		return nil, err
	}

	if len(searchResults) == 0 {
		return []provider.SeasonMetadata{}, nil
	}

	// TODO: map result

	return nil, nil
}

func (p *TVMazeProvider) SearchShow(name string, year int) ([]provider.ShowSearchResult, error) {
	params := url.Values{}
	params.Add("q", name)
	encoded := params.Encode()

	fullUrl := "https://api.tvmaze.com/search/shows?" + encoded

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("tvmaze search returned status %d", resp.StatusCode)
		return nil, err
	}

	var searchResults []TVMazeSearchResult
	if err = json.Unmarshal(body, &searchResults); err != nil {
		return nil, err
	}

	if len(searchResults) == 0 {
		return []provider.ShowSearchResult{}, nil
	}

	mappedResults := make([]provider.ShowSearchResult, 0, len(searchResults))
	for _, result := range searchResults {

		// if a year is provided -> check the year of the fetched show and continue if it doesn't match
		if year != 0 && result.Show.Premiered != nil {
			resultYear, err := getResultYear(*result.Show.Premiered)
			logger.Debug(nil, "Result Year: {Resultyear}", resultYear)
			if err == nil {
				if resultYear != year {
					continue
				}
			} else {
				logger.Error(err, "Couldn't get result year")
			}
		}

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

	return mappedResults, nil
}

func getResultYear(premiered string) (int, error) {
	split := strings.Split(premiered, "-")
	if len(split) > 0 {
		resultYear, err := strconv.Atoi(split[0])
		if err != nil {
			return 0, err
		}
		return resultYear, nil
	} else {
		return 0, fmt.Errorf("String is in wrong format.")
	}
}

var _ provider.TVMetadataProvider = (*TVMazeProvider)(nil)
