package subdl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/helper"
	"github.com/loissascha/localstream/internal/provider"
	"github.com/loissascha/localstream/internal/repository"
)

type SubDlProvider struct {
	apiKey            string
	movieSubtitleRepo repository.MovieSubtitleRepository
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

func NewSubDlProvider(
	apiKey string,
	movieSubtitleRepo repository.MovieSubtitleRepository,
) *SubDlProvider {
	return &SubDlProvider{
		apiKey:            apiKey,
		movieSubtitleRepo: movieSubtitleRepo,
	}
}

func (self *SubDlProvider) DownloadMovieSubtitle(ctx context.Context, movieId uuid.UUID, providerResult provider.SubtitleProviderResult) error {
	fullUrl := "https://dl.subdl.com" + providerResult.Url
	movieIdStr := encoders.EncodeUUID(movieId)

	downloadedPath, err := helper.DownloadSubtitleForMovie(fullUrl, movieIdStr)
	if err != nil {
		return err
	}
	logger.Info(nil, "Downloaded to {DownloadedPath}", downloadedPath)

	ext, err := helper.GetExtensionFromUrl(downloadedPath)
	if err != nil {
		return err
	}
	ext = strings.TrimLeft(ext, ".")

	if ext == "zip" {
		err := self.processZip(ctx, downloadedPath, movieId, providerResult)
		if err != nil {
			return err
		}
	} else if ext == "srt" {
		err := self.processSrt(ctx, downloadedPath, movieId, providerResult)
		if err != nil {
			return err
		}
	} else if ext == "vtt" {
		err := self.processVtt(ctx, downloadedPath, movieId, providerResult)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *SubDlProvider) processVtt(ctx context.Context, downloadedPath string, movieId uuid.UUID, providerResult provider.SubtitleProviderResult) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	downloadedPath = strings.TrimLeft(downloadedPath, "/")
	downloadedPath = "/" + downloadedPath

	sub, err := self.movieSubtitleRepo.GetByPath(ctx, downloadedPath)
	if err != nil {
		return err
	}
	if sub != nil {
		return nil
	}

	subt := entity.MovieSubtitle{
		ID:        id,
		MovieID:   movieId,
		Path:      downloadedPath,
		Name:      providerResult.Name,
		Lang:      providerResult.Lang,
		LangShort: providerResult.LangShort,
	}
	err = self.movieSubtitleRepo.Create(ctx, &subt)
	if err != nil {
		return err
	}

	return nil
}

func (self *SubDlProvider) processSrt(ctx context.Context, downloadedPath string, movieId uuid.UUID, providerResult provider.SubtitleProviderResult) error {
	newPath, err := self.convertSubtitleSrt(downloadedPath)
	if err != nil {
		return err
	}
	logger.Info(nil, "New Path after srt convert: {NewPath}", newPath)
	ext, err := helper.GetExtensionFromUrl(newPath)
	if err != nil {
		return err
	}
	ext = strings.TrimLeft(ext, ".")
	if ext == "vtt" {
		self.processVtt(ctx, newPath, movieId, providerResult)
	} else {
		logger.Error(nil, "Resulting filea fter srt to vtt convert is not vtt...")
	}
	return nil
}

func (self *SubDlProvider) processZip(ctx context.Context, downloadedPath string, movieId uuid.UUID, providerResult provider.SubtitleProviderResult) error {
	paths, err := self.unpackSubtitleZip(downloadedPath)
	if err != nil {
		return err
	}
	for _, p := range paths {
		ext, err := helper.GetExtensionFromUrl(p)
		if err != nil {
			return err
		}
		ext = strings.TrimLeft(ext, ".")
		if ext == "srt" {
			self.processSrt(ctx, p, movieId, providerResult)
		} else if ext == "vtt" {
			self.processVtt(ctx, p, movieId, providerResult)
		} else {
			logger.Error(nil, "Resulting file after unzip is neither srt nor vtt but {ext}", ext)
		}
	}
	return nil
}

func (self *SubDlProvider) unpackSubtitleZip(localPath string) ([]string, error) {
	p, err := helper.UnzipMultiFiles(localPath)
	return p, err
}

func (self *SubDlProvider) convertSubtitleSrt(localPath string) (string, error) {
	p, err := helper.ConvertToVTT(localPath)
	if err != nil {
		return "", err
	}
	return p, nil
}

func (self *SubDlProvider) SearchEpisode(ctx context.Context, showName string, seasonNumber int, episodeNumber int, lang string) ([]provider.SubtitleProviderResult, error) {
	params := url.Values{}
	params.Add("api_key", self.apiKey)
	params.Add("film_name", showName)
	if seasonNumber > 0 {
		params.Add("season_number", strconv.Itoa(seasonNumber))
		if episodeNumber > 0 {
			params.Add("episode_number", strconv.Itoa(episodeNumber))
		}
	}
	params.Add("type", "tv")
	params.Add("languages", lang)
	encoded := params.Encode()

	fullUrl := "https://api.subdl.com/api/v1/subtitles?" + encoded
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

	var results = []provider.SubtitleProviderResult{}
	for _, s := range searchResults.Subtitles {
		results = append(results, provider.SubtitleProviderResult{
			Name:      s.ReleaseName,
			Lang:      s.Lang,
			LangShort: s.Language,
			Author:    s.Author,
			Url:       s.Url,
		})
	}

	return results, nil
}

func (self *SubDlProvider) SearchMovie(ctx context.Context, name string, lang string) ([]provider.SubtitleProviderResult, error) {
	params := url.Values{}
	params.Add("api_key", self.apiKey)
	params.Add("film_name", name)
	params.Add("type", "movie")
	params.Add("languages", lang)
	encoded := params.Encode()

	fullUrl := "https://api.subdl.com/api/v1/subtitles?" + encoded

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

	var results = []provider.SubtitleProviderResult{}
	for _, s := range searchResults.Subtitles {
		results = append(results, provider.SubtitleProviderResult{
			Name:      s.ReleaseName,
			Lang:      s.Lang,
			LangShort: s.Language,
			Author:    s.Author,
			Url:       s.Url,
		})
	}

	return results, nil
}

var _ provider.SubtitleProvider = (*SubDlProvider)(nil)
