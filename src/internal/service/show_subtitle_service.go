package service

import (
	"context"
	"strconv"

	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/provider"
	"github.com/loissascha/localstream/internal/repository"
)

type ShowSubtitleService struct {
	subtitleProvider provider.SubtitleProvider
	episodeRepo      repository.EpisodeRepository
	seasonRepo       repository.SeasonRepository
}

func NewShowSubtitleService(
	subtitleProvider provider.SubtitleProvider,
	episodeRepo repository.EpisodeRepository,
	seasonRepo repository.SeasonRepository,
) *ShowSubtitleService {
	return &ShowSubtitleService{
		subtitleProvider: subtitleProvider,
		episodeRepo:      episodeRepo,
		seasonRepo:       seasonRepo,
	}
}

func (s *ShowSubtitleService) CreateFromSubtitleResult(ctx context.Context, episodeID string, subtitle provider.SubtitleProviderResult) error {
	// showUUID, err := encoders.DecodeUUID(showID)
	// if err != nil {
	// 	return err
	// }

	episodeUUID, err := encoders.DecodeUUID(episodeID)
	if err != nil {
		return err
	}

	episode, err := s.episodeRepo.GetByID(ctx, episodeUUID)
	if err != nil {
		return err
	}

	seasonID := episode.SeasonID
	season, err := s.seasonRepo.GetByID(ctx, seasonID)
	if err != nil {
		return err
	}

	showID := season.ShowID

	err = s.subtitleProvider.DownloadShowSubtitle(ctx, showID, season, episode, subtitle)
	if err != nil {
		return err
	}

	return nil
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
