package service

import (
	"context"
	"strconv"

	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/provider"
)

type ShowSubtitleService struct {
	subtitleProvider provider.SubtitleProvider
}

func NewShowSubtitleService(
	subtitleProvider provider.SubtitleProvider,
) *ShowSubtitleService {
	return &ShowSubtitleService{
		subtitleProvider: subtitleProvider,
	}
}

func (s *ShowSubtitleService) SearchByEpisode(ctx context.Context, show string, season string, episode string, lang string) ([]provider.SubtitleProviderResult, error) {
	seasonInt, err := strconv.Atoi(season)
	if err != nil {
		logger.Warning(err, "Error converting season to int from {season}", season)
		seasonInt = 0
	}
	episodeInt, err := strconv.Atoi(episode)
	if err != nil {
		logger.Warning(err, "Error converting episode to int from {episode}", episode)
		episodeInt = 0
	}
	return s.subtitleProvider.SearchEpisode(ctx, show, seasonInt, episodeInt, lang)
}
