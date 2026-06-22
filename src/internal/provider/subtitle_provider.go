package provider

import (
	"context"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type SubtitleProviderResult struct {
	Name      string `json:"name"`
	Name2     string `json:"name2"`
	Lang      string `json:"lang"`
	LangShort string `json:"lang_short"`
	Author    string `json:"author"`
	Url       string `json:"url"`
}

type SubtitleSupportedLanguage struct {
	Value string `json:"value"`
	Name  string `json:"name"`
}

type SubtitleProvider interface {
	SearchMovie(ctx context.Context, name string, lang string) ([]SubtitleProviderResult, error)
	SearchEpisode(ctx context.Context, showName string, seasonNumber int, episodeNumber int, lang string) ([]SubtitleProviderResult, error)
	DownloadMovieSubtitle(ctx context.Context, movieId uuid.UUID, providerResult SubtitleProviderResult) error
	DownloadShowSubtitle(ctx context.Context, showId uuid.UUID, season *entity.Season, episode *repository.EpisodeWithMetadata, providerResult SubtitleProviderResult) error
	SupportedLanguages(ctx context.Context) ([]SubtitleSupportedLanguage, error)
}
