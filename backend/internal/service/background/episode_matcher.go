package background

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/provider"
	"github.com/loissascha/localstream/internal/repository"
)

type episodeMetadataCache struct {
	created  time.Time
	metadata []provider.EpisodeMetadata
}

type EpisodeMatcher struct {
	Channel              chan *entity.Episode
	episodeMetadataCache map[int]episodeMetadataCache
	metadataProvider     provider.TVMetadataProvider
	seasonRepo           repository.SeasonRepository
	seasonMetadataRepo   repository.SeasonMetadataRepository
	showRepo             repository.ShowRepository
	showMetadataRepo     repository.ShowMetadataRepository
	episodeRepo          repository.EpisodeRepository
	episodeMetadataRepo  repository.EpisodeMetadataRepository
}

func NewEpisodeMatcher(metadataProvider provider.TVMetadataProvider, seasonMetaRepo repository.SeasonMetadataRepository, seasonRepo repository.SeasonRepository, showRepo repository.ShowRepository, showMetaRepo repository.ShowMetadataRepository, episodeRepo repository.EpisodeRepository, episodeMetadataRepo repository.EpisodeMetadataRepository) *EpisodeMatcher {
	ch := make(chan *entity.Episode)
	return &EpisodeMatcher{
		Channel:              ch,
		episodeMetadataCache: make(map[int]episodeMetadataCache),
		metadataProvider:     metadataProvider,
		seasonRepo:           seasonRepo,
		seasonMetadataRepo:   seasonMetaRepo,
		showRepo:             showRepo,
		showMetadataRepo:     showMetaRepo,
		episodeRepo:          episodeRepo,
		episodeMetadataRepo:  episodeMetadataRepo,
	}
}

func (self *EpisodeMatcher) getMetadataResultLive(fetchID int) ([]provider.EpisodeMetadata, error) {
	result, err := self.metadataProvider.SearchEpisodes(fetchID)
	self.episodeMetadataCache[fetchID] = episodeMetadataCache{
		created:  time.Now().UTC(),
		metadata: result,
	}
	return result, err
}

func (self *EpisodeMatcher) getMetadataResultCacheOrLive(fetchID int) ([]provider.EpisodeMetadata, error) {
	cachefile, ok := self.episodeMetadataCache[fetchID]
	if ok {
		if time.Now().UTC().Sub(cachefile.created) > 24*time.Hour {
			return self.getMetadataResultLive(fetchID)
		}
		logger.Debug(nil, "Load season metadata from cache")
		return cachefile.metadata, nil
	}

	return self.getMetadataResultLive(fetchID)
}

func (self *EpisodeMatcher) RunBackground() {
	go func() {
		for {
			episode := <-self.Channel
			if !episode.FetchSource.IsNone() {
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			season, err := self.seasonRepo.GetByID(ctx, episode.SeasonID)
			if err != nil {
				logger.Error(err, "Can't find season")
				continue
			}

			show, err := self.showRepo.GetByID(ctx, season.ShowID)
			if err != nil {
				logger.Error(err, "Can't find show")
				continue
			}

			// show has multiple fetch sources -> wait until it's clear before loading any metadata
			if show.FetchSource.IsEmpty() || show.FetchSource.IsMultiple() || show.FetchSource.IsNone() {
				continue
			}

			seasonMetadata, err := self.seasonMetadataRepo.GetBySeasonID(ctx, season.ID)
			if err != nil || seasonMetadata == nil {
				logger.Error(err, "Can't get season metadata")
				continue
			}

			episodeMetadataResult, err := self.getMetadataResultCacheOrLive(seasonMetadata.FetchID)
			if err != nil {
				logger.Error(err, "Can't get episode metadata")
				continue
			}

			hasError := false
			for _, emr := range episodeMetadataResult {
				if emr.Number == episode.Number {
					err := self.createEpisodeMetadata(ctx, episode, &emr)
					if err != nil {
						logger.Error(err, "Error creating metadata for episode")
						hasError = true
					}
					break
				}
			}

			if !hasError {
				episode.FetchSource = entity.FetchSourceTVMaze
				err := self.episodeRepo.UpdateFetchSource(ctx, episode.ID, episode.FetchSource)
				if err != nil {
					logger.Error(err, "Error updating fetch source of episode")
					continue
				}
			}

		}
	}()
}

func (self *EpisodeMatcher) createEpisodeMetadata(ctx context.Context, episode *entity.Episode, metadata *provider.EpisodeMetadata) error {
	mid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	mediumImage := ""
	originalImage := ""
	if metadata.Image != nil {
		mediumImage = metadata.Image.Medium
		originalImage = metadata.Image.Original
	}
	m := entity.EpisodeMetadata{
		ID:               mid,
		EpisodeID:        episode.ID,
		Url:              metadata.Url,
		Number:           metadata.Number,
		Summary:          metadata.Summary,
		MediumImageUrl:   mediumImage,
		OriginalImageUrl: originalImage,
		FetchSource:      entity.FetchSourceTVMaze,
		FetchID:          metadata.ID,
	}
	err = self.episodeMetadataRepo.Create(ctx, &m)
	if err != nil {
		return err
	}
	return nil
}
