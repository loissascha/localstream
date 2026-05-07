package background

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/provider"
	"github.com/loissascha/localstream/internal/repository"
)

type seasonMetadataCache struct {
	created  time.Time
	metadata []provider.SeasonMetadata
}

type SeasonMatcher struct {
	Channel             chan *entity.Season
	seasonMetadataCache map[int]seasonMetadataCache
	metadataProvider    provider.TVMetadataProvider
	seasonRepo          repository.SeasonRepository
	seasonMetadataRepo  repository.SeasonMetadataRepository
	showRepo            repository.ShowRepository
	showMetadataRepo    repository.ShowMetadataRepository
}

func NewSeasonMatcher(metadataProvider provider.TVMetadataProvider, seasonMetaRepo repository.SeasonMetadataRepository, seasonRepo repository.SeasonRepository, showRepo repository.ShowRepository, showMetaRepo repository.ShowMetadataRepository) *SeasonMatcher {
	ch := make(chan *entity.Season)
	return &SeasonMatcher{
		Channel:             ch,
		seasonMetadataCache: make(map[int]seasonMetadataCache),
		metadataProvider:    metadataProvider,
		seasonRepo:          seasonRepo,
		seasonMetadataRepo:  seasonMetaRepo,
		showRepo:            showRepo,
		showMetadataRepo:    showMetaRepo,
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
					err := self.createSeasonMetadata(ctx, season, &smr)
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

func (self *SeasonMatcher) createSeasonMetadata(ctx context.Context, season *entity.Season, metadata *provider.SeasonMetadata) error {
	mid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	mediumImage := ""
	originalImage := ""
	if metadata.Image != nil {
		mediumImage, err = downloadImageAndGetStaticPath(metadata.Image.Medium, fmt.Sprintf("med_SE_%s", mid.String()))
		if err != nil {
			return err
		}
		originalImage, err = downloadImageAndGetStaticPath(metadata.Image.Original, fmt.Sprintf("org_%s", mid.String()))
		if err != nil {
			return err
		}
	}
	m := entity.SeasonMetadata{
		ID:               mid,
		SeasonID:         season.ID,
		Url:              metadata.Url,
		Number:           metadata.Number,
		Summary:          metadata.Summary,
		PremiereDate:     metadata.PremiereDate,
		MediumImageUrl:   mediumImage,
		OriginalImageUrl: originalImage,
		FetchSource:      entity.FetchSourceTVMaze,
		FetchID:          metadata.ID,
	}
	err = self.seasonMetadataRepo.Create(ctx, &m)
	if err != nil {
		return err
	}
	return nil
}
