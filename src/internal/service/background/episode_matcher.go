package background

import (
	"context"
	"time"

	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/provider"
	"github.com/loissascha/localstream/internal/repository"
	"github.com/loissascha/localstream/internal/service"
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
	episodeMetaService   *service.EpisodeMetadataService
}

func NewEpisodeMatcher(
	metadataProvider provider.TVMetadataProvider,
	seasonMetaRepo repository.SeasonMetadataRepository,
	seasonRepo repository.SeasonRepository,
	showRepo repository.ShowRepository,
	showMetaRepo repository.ShowMetadataRepository,
	episodeRepo repository.EpisodeRepository,
	episodeMetadataRepo repository.EpisodeMetadataRepository,
	episodeMetaService *service.EpisodeMetadataService,
) *EpisodeMatcher {
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
		episodeMetaService:   episodeMetaService,
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
		// logger.Debug(nil, "Load episode metadata from cache")
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
				logger.Error(err, "[EpisodeMatcher] Can't find season")
				continue
			}

			show, err := self.showRepo.GetByID(ctx, season.ShowID)
			if err != nil {
				logger.Error(err, "[EpisodeMatcher] Can't find show")
				continue
			}

			// show | season has multiple fetch sources -> wait until it's clear before loading any metadata
			if show.FetchSource.IsEmpty() || show.FetchSource.IsMultiple() || show.FetchSource.IsNone() {
				continue
			}
			if season.FetchSource.IsEmpty() || season.FetchSource.IsNone() || season.FetchSource.IsMultiple() {
				continue
			}

			seasonMetadata, err := self.seasonMetadataRepo.GetBySeasonID(ctx, season.ID)
			if err != nil {
				logger.Error(err, "[EpisodeMatcher] Can't get season metadata")
				continue
			}
			if seasonMetadata == nil {
				err := self.seasonRepo.UpdateFetchSource(ctx, season.ID, entity.FetchSourceNone)
				if err != nil {
					logger.Error(err, "[EpisodeMatcher] Erorr resetting Season Metadata")
					continue
				}
				logger.Warning(nil, "[EpisodeMatcher] Season Metadata for season {SeasonNumber} of show {ShowName} with fetchsource: {FetchSoruce} was not found. Resetting Season fetch source.", season.Number, show.Name, season.FetchSource)
				continue
			}

			episodeMetadataResult, err := self.getMetadataResultCacheOrLive(seasonMetadata.FetchID)
			if err != nil {
				logger.Error(err, "[EpisodeMatcher] Can't get episode metadata")
				continue
			}

			hasError := false
			for _, emr := range episodeMetadataResult {
				if emr.Number == episode.Number {
					err := self.episodeMetaService.CreateEpisodeMetadata(ctx, episode, &emr)
					if err != nil {
						logger.Error(err, "[EpisodeMatcher] Error creating metadata for episode")
						hasError = true
					}
					break
				}
			}

			if !hasError {
				episode.FetchSource = entity.FetchSourceTVMaze
				err := self.episodeRepo.UpdateFetchSource(ctx, episode.ID, episode.FetchSource)
				if err != nil {
					logger.Error(err, "[EpisodeMatcher] Error updating fetch source of episode")
					continue
				}
			}
		}
	}()
}
