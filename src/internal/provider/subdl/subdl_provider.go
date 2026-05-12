package subdl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/loissascha/localstream/internal/provider"
)

type SubDlProvider struct {
	apiKey string
}

type ApiSearchResult struct {
	Status    bool                `json:"status"`
	Subtitles []ApiSubtitleResult `json:"subtitles"`
}

type ApiSubtitleResult struct {
	ReleaseName  string `json:"release_name"`
	Name         string `json:"name"`
	Lang         string `json:"lang"`
	Language     string `json:"language"`
	Author       string `json:"author"`
	Url          string `json:"url"`
	SubtitlePage string `json:"subtitlePage"`
}

func NewSubDlProvider(apiKey string) *SubDlProvider {
	return &SubDlProvider{
		apiKey: apiKey,
	}
}

func (self *SubDlProvider) SearchMovie(name string) ([]provider.SubtitleProviderResult, error) {
	params := url.Values{}
	params.Add("api_key", self.apiKey)
	params.Add("film_name", name)
	params.Add("type", "movie")
	params.Add("languages", "EN")
	encoded := params.Encode()

	fullUrl := "https://api.subdl.com/api/v1/subtitles?" + encoded

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
		err = fmt.Errorf("subdl search returned status %d", resp.StatusCode)
		return nil, err
	}

	var searchResults ApiSearchResult
	if err = json.Unmarshal(body, &searchResults); err != nil {
		return nil, err
	}
	if !searchResults.Status {
		return nil, fmt.Errorf("api call failed: status false")
	}

	var results = []provider.SubtitleProviderResult{}
	for _, s := range searchResults.Subtitles {
		results = append(results, provider.SubtitleProviderResult{
			Name:   s.ReleaseName,
			Lang:   s.Lang,
			Author: s.Author,
			Url:    s.Url,
		})
	}

	return results, nil
}

var _ provider.SubtitleProvider = (*SubDlProvider)(nil)
