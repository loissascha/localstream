package background

import (
	"context"

	"github.com/google/uuid"
	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/parsers"
)

func (s *BackgroundService) createEpisode(ctx context.Context, episodeInfo *parsers.EpisodeInfo, seasonId uuid.UUID, episodePath string) (uuid.UUID, error) {

	episode := &entity.Episode{
		ID:          uuid.New(),
		SeasonID:    seasonId,
		Number:      episodeInfo.Episode,
		Path:        episodePath,
		FetchSource: entity.FetchSourceNone,
	}

	err := s.episodeRepo.Create(ctx, episode)
	if err != nil {
		logger.Error(err, "Error creating episode")
		return uuid.Nil, err
	}
	return episode.ID, nil
}

func (s *BackgroundService) createSeason(ctx context.Context, seasonInfo *parsers.SeasonInfo, showId uuid.UUID, seasonPath string) (uuid.UUID, error) {
	season := &entity.Season{
		ID:          uuid.New(),
		ShowID:      showId,
		Number:      seasonInfo.Season,
		Path:        seasonPath,
		FetchSource: entity.FetchSourceNone,
	}

	err := s.seasonRepo.Create(ctx, season)
	if err != nil {
		logger.Error(err, "Error creating season")
		return uuid.Nil, err
	}
	return season.ID, nil
}

func (s *BackgroundService) createShow(ctx context.Context, showInfo *parsers.ShowInfo, showPath string) (uuid.UUID, error) {
	showE := &entity.Show{
		ID:          uuid.New(),
		Name:        showInfo.Series,
		Year:        0,
		Path:        showPath,
		FetchSource: entity.FetchSourceNone,
	}

	if showInfo.Year != nil {
		showE.Year = *showInfo.Year
	}

	err := s.showRepo.Create(ctx, showE)
	if err != nil {
		logger.Error(err, "Error creating show")
		return uuid.Nil, err
	}
	return showE.ID, nil
}

func (s *BackgroundService) findShowWithPathInList(path string, list []entity.Show) (bool, uuid.UUID) {
	for _, show := range list {
		if show.Path == path {
			return true, show.ID
		}
	}
	return false, uuid.Nil
}

func (s *BackgroundService) findSeasonWithPathInList(path string, list []entity.Season) (bool, uuid.UUID) {
	for _, season := range list {
		if season.Path == path {
			return true, season.ID
		}
	}
	return false, uuid.Nil
}

func (s *BackgroundService) findEpisodeWithPathInList(path string, list []entity.Episode) (bool, uuid.UUID) {
	for _, e := range list {
		if e.Path == path {
			return true, e.ID
		}
	}
	return false, uuid.Nil
}

func (s *BackgroundService) findMovieWithPathInList(path string, list []entity.Movie) (bool, uuid.UUID) {
	for _, m := range list {
		if m.Path == path {
			return true, m.ID
		}
	}
	return false, uuid.Nil
}
