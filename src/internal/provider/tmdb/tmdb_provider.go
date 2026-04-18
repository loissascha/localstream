package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/provider"
)

type MovieAPISearchResult struct {
	Page    int              `json:"page"`
	Results []MovieAPIResult `json:"results"`
}

type MovieAPIResult struct {
	ID            int    `json:"id"`
	Adult         bool   `json:"adult"`
	OriginalTitle string `json:"original_title"`
	Overview      string `json:"overview"`
	ReleaseDate   string `json:"release_date"`
	BackdropPath  string `json:"backdrop_path"`
	PosterPath    string `json:"poster_path"`
}

type TMDBProvider struct {
}

func NewTMDBProvider() *TMDBProvider {
	return &TMDBProvider{}
}

func (self *TMDBProvider) GetMovieByID(id int) (*provider.MovieResult, error) {
	key := os.Getenv("TMDB_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("No TMDB api key configured.")
	}
	params := url.Values{}
	params.Add("api_key", key)
	encoded := params.Encode()
	fullUrl := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d?%s", id, encoded)

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

	var searchResult MovieAPIResult
	if err = json.Unmarshal(body, &searchResult); err != nil {
		return nil, err
	}
	movieResult := toMovieResult(searchResult)

	return &movieResult, nil
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

	var searchResults MovieAPISearchResult
	if err = json.Unmarshal(body, &searchResults); err != nil {
		return nil, err
	}

	var results = []provider.MovieResult{}
	for _, sr := range searchResults.Results {
		results = append(results, toMovieResult(sr))
	}

	if len(searchResults.Results) == 0 {
		return results, nil
	}

	return results, nil
}

func toMovieResult(r MovieAPIResult) provider.MovieResult {
	ry := 0
	rd := r.ReleaseDate
	sp := strings.Split(rd, "-")
	var err error
	if len(sp) == 3 {
		ry, err = strconv.Atoi(sp[0])
		if err != nil {
			logger.Error(err, "Error parsing release year for movie")
		}
	}
	return provider.MovieResult{
		FetchID:      r.ID,
		Adult:        r.Adult,
		Title:        r.OriginalTitle,
		Description:  r.Overview,
		ReleaseYear:  ry,
		BackdropPath: r.BackdropPath,
		PosterPath:   r.PosterPath,
	}
}

var _ provider.MovieMetadataProvider = (*TMDBProvider)(nil)
