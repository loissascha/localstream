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

type SeasonMatcher struct {
	Channel            chan *entity.Season
	metadataProvider   provider.TVMetadataProvider
	seasonRepo         repository.SeasonRepository
	seasonMetadataRepo repository.SeasonMetadataRepository
	showRepo           repository.ShowRepository
	showMetadataRepo   repository.ShowMetadataRepository
}

func NewSeasonMatcher(metadataProvider provider.TVMetadataProvider, seasonMetaRepo repository.SeasonMetadataRepository, seasonRepo repository.SeasonRepository, showRepo repository.ShowRepository, showMetaRepo repository.ShowMetadataRepository) *SeasonMatcher {
	ch := make(chan *entity.Season)
	return &SeasonMatcher{
		Channel:            ch,
		metadataProvider:   metadataProvider,
		seasonRepo:         seasonRepo,
		seasonMetadataRepo: seasonMetaRepo,
		showRepo:           showRepo,
		showMetadataRepo:   showMetaRepo,
	}
}

func (self *SeasonMatcher) RunBackground() {
	go func() {
		for {
			season := <-self.Channel

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			show, err := self.showRepo.GetByID(ctx, season.ShowID)
			if err != nil {
				logger.Error(err, "Can't find show")
				continue
			}

			// show has multiple fetch sources -> wait until it's clear before loading any metadata
			if show.FetchSource.IsEmpty() || show.FetchSource.IsMultiple() || show.FetchSource.IsNone() {
				continue
			}

			showMetadata, err := self.showMetadataRepo.GetByShowID(ctx, show.ID)
			if err != nil {
				logger.Error(err, "Can't get show metadata")
				continue
			}
			if len(showMetadata) != 1 {
				logger.Error(nil, "Show has wrong amount of metadatas: {Len}", len(showMetadata))
				continue
			}

			seasonMetadataResult, err := self.metadataProvider.SearchSeasons(showMetadata[0].FetchID)
			if err != nil {
				logger.Error(err, "Can't get season metadatas")
				continue
			}

			for _, smr := range seasonMetadataResult {
				if smr.Number == season.Number {
					err := self.createSeasonMetadata(ctx, season, &smr)
					if err != nil {
						logger.Error(err, "Error creating metadata for season")
					}
					break
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
		mediumImage = metadata.Image.Medium
		originalImage = metadata.Image.Original
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
