package background

import (
	"context"
	"time"

	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/entity"
)

func (s *BackgroundService) runEpisodesMatcher() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	episodes, err := s.episodeRepo.NecessaryForMetadataFetch(ctx)
	if err != nil {
		return err
	}

	for _, episode := range episodes {
		err := s.fetchMetadataForEpisode(ctx, &episode)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BackgroundService) fetchMetadataForEpisode(ctx context.Context, episode *entity.Episode) error {
	if !episode.FetchSource.IsNone() {
		return nil
	}

	season, err := s.seasonRepo.GetByID(ctx, episode.SeasonID)
	if err != nil {
		return err
	}

	show, err := s.showRepo.GetByID(ctx, season.ShowID)
	if err != nil {
		return err
	}

	// show | season has multiple fetch sources -> wait until it's clear before loading any metadata
	if !show.FetchSource.IsValid() {
		return nil
	}
	if !season.FetchSource.IsValid() {
		return nil
	}

	seasonMetadata, err := s.seasonMetadataRepo.GetBySeasonID(ctx, season.ID)
	if err != nil {
		return err
	}
	if seasonMetadata == nil {
		err := s.seasonRepo.UpdateFetchSource(ctx, season.ID, entity.FetchSourceNone)
		if err != nil {
			return err
		}
		logger.Warning(nil, "[EpisodeMatcher] Season Metadata for season {SeasonNumber} of show {ShowName} with fetchsource: {FetchSoruce} was not found. Resetting Season fetch source.", season.Number, show.Name, season.FetchSource)
		return nil
	}

	episodeMetadataResult, err := s.getEpisodeMetadataResultCacheOrLive(seasonMetadata.FetchID)
	if err != nil {
		return err
	}

	for _, emr := range episodeMetadataResult {
		if emr.Number == episode.Number {
			err := s.episodeMetadataService.CreateEpisodeMetadata(ctx, episode, &emr)
			if err != nil {
				return err
			}
			break
		}
	}

	episode.FetchSource = entity.FetchSourceTVMaze
	err = s.episodeRepo.UpdateFetchSource(ctx, episode.ID, episode.FetchSource)
	if err != nil {
		return err
	}
	return nil
}
