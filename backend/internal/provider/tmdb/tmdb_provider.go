package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/loissascha/localstream/internal/provider"
)

type TMDBProvider struct {
}

func NewTMDBProvider() *TMDBProvider {
	return &TMDBProvider{}
}

func (self *TMDBProvider) SearchMovie(name string, year int) ([]provider.MovieResult, error) {
	key := os.Getenv("TMDB_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("No TMDB api key configured.")
	}
	params := url.Values{}
	params.Add("api_key", key)
	params.Add("query", name)
	params.Add("year", fmt.Sprintf("%d", year))
	encoded := params.Encode()

	fullUrl := "https://api.themoviedb.org/3/search/movie?" + encoded

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

	var searchResults provider.MovieSearchResult
	if err = json.Unmarshal(body, &searchResults); err != nil {
		return nil, err
	}

	if len(searchResults.Results) == 0 {
		return searchResults.Results, nil
	}

	return searchResults.Results, nil
}

var _ provider.MovieMetadataProvider = (*TMDBProvider)(nil)
