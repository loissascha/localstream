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

type seasonMetadataCache struct {
	created  time.Time
	metadata []provider.SeasonMetadata
}

type SeasonMatcher struct {
	Channel               chan *entity.Season
	seasonMetadataCache   map[int]seasonMetadataCache
	metadataProvider      provider.TVMetadataProvider
	seasonRepo            repository.SeasonRepository
	seasonMetadataRepo    repository.SeasonMetadataRepository
	showRepo              repository.ShowRepository
	showMetadataRepo      repository.ShowMetadataRepository
	seasonMetadataService *service.SeasonMetadataService
}

func NewSeasonMatcher(
	metadataProvider provider.TVMetadataProvider,
	seasonMetaRepo repository.SeasonMetadataRepository,
	seasonRepo repository.SeasonRepository,
	showRepo repository.ShowRepository,
	showMetaRepo repository.ShowMetadataRepository,
	seasonMetaService *service.SeasonMetadataService,
) *SeasonMatcher {
	ch := make(chan *entity.Season)
	return &SeasonMatcher{
		Channel:               ch,
		seasonMetadataCache:   make(map[int]seasonMetadataCache),
		metadataProvider:      metadataProvider,
		seasonRepo:            seasonRepo,
		seasonMetadataRepo:    seasonMetaRepo,
		showRepo:              showRepo,
		showMetadataRepo:      showMetaRepo,
		seasonMetadataService: seasonMetaService,
	}
}

func (self *SeasonMatcher) getMetadataResultLive(fetchID int) ([]provider.SeasonMetadata, error) {
	seasonMetadataResult, err := self.metadataProvider.SearchSeasons(fetchID)
	self.seasonMetadataCache[fetchID] = seasonMetadataCache{
		created:  time.Now().UTC(),
		metadata: seasonMetadataResult,
	}
	return seasonMetadataResult, err
}

func (self *SeasonMatcher) getMetadataResultCacheOrLive(fetchID int) ([]provider.SeasonMetadata, error) {
	cachefile, ok := self.seasonMetadataCache[fetchID]
	if ok {
		if time.Now().UTC().Sub(cachefile.created) > 24*time.Hour {
			return self.getMetadataResultLive(fetchID)
		}
		return cachefile.metadata, nil
	}

	return self.getMetadataResultLive(fetchID)
}

func (self *SeasonMatcher) RunBackground() {
	go func() {
		for {
			season := <-self.Channel
			if !season.FetchSource.IsNone() {
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			show, err := self.showRepo.GetByID(ctx, season.ShowID)
			if err != nil {
				logger.Error(err, "[SeasonMatcher] Can't find show")
				continue
			}

			// show has multiple fetch sources -> wait until it's clear before loading any metadata
			if show.FetchSource.IsEmpty() || show.FetchSource.IsMultiple() || show.FetchSource.IsNone() {
				continue
			}

			showMetadata, err := self.showMetadataRepo.GetByShowID(ctx, show.ID)
			if err != nil {
				logger.Error(err, "[SeasonMatcher] Can't get show metadata")
				continue
			}
			if len(showMetadata) != 1 {
				logger.Error(nil, "[SeasonMatcher] Show has wrong amount of metadatas: {Len}", len(showMetadata))
				continue
			}

			seasonMetadataResult, err := self.getMetadataResultCacheOrLive(showMetadata[0].FetchID)
			if err != nil {
				logger.Error(err, "[SeasonMatcher] Can't get season metadatas")
				continue
			}

			hasError := false
			for _, smr := range seasonMetadataResult {
				if smr.Number == season.Number {
					err := self.seasonMetadataService.CreateSeasonMetadata(ctx, season, &smr)
					if err != nil {
						logger.Error(err, "[SeasonMatcher] Error creating metadata for season")
						hasError = true
					}
					break
				}
			}

			if !hasError {
				season.FetchSource = entity.FetchSourceTVMaze
				err := self.seasonRepo.UpdateFetchSource(ctx, season.ID, season.FetchSource)
				if err != nil {
					logger.Error(err, "[SeasonMatcher] Error updating fetch source of season")
					continue
				}
			}
		}
	}()
}
