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

type ShowMatcher struct {
	Channel            chan *entity.Show
	metadataProvider   provider.TVMetadataProvider
	showRepo           repository.ShowRepository
	showMetadataRepo   repository.ShowMetadataRepository
	seasonMetadataRepo repository.SeasonMetadataRepository
}

func NewShowMatcher(metadataProvider provider.TVMetadataProvider, showRepo repository.ShowRepository, showMetadataRepo repository.ShowMetadataRepository, seasonMetaRepo repository.SeasonMetadataRepository) *ShowMatcher {
	ch := make(chan *entity.Show)
	return &ShowMatcher{
		Channel:            ch,
		metadataProvider:   metadataProvider,
		showRepo:           showRepo,
		showMetadataRepo:   showMetadataRepo,
		seasonMetadataRepo: seasonMetaRepo,
	}
}

func (self *ShowMatcher) RunBackground() {
	go func() {
		for {
			show := <-self.Channel
			logger.Info(nil, "New ShowID triggered! {Show}", show)

			showSearchResults, err := self.metadataProvider.SearchShow(show.Name, show.Year)
			if err != nil {
				logger.Error(err, "Error getting show results")
				continue
			}

			if len(showSearchResults) == 0 {
				logger.Error(nil, "Didn't find anything for show {Show} ({Year})", show.Name, show.Year)
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			hasError := false
			createdMetadata := []*entity.ShowMetadata{}
			for _, res := range showSearchResults {
				m, err := self.createShowMetadata(ctx, show, res)
				if err != nil {
					logger.Error(err, "Error creating show metadata")
					hasError = true
				}
				createdMetadata = append(createdMetadata, m)
			}

			if hasError {
				continue
			}

			if len(showSearchResults) > 1 {
				show.FetchSource = entity.FetchSourceMultiple
				self.showRepo.UpdateFetchSource(ctx, show.ID, entity.FetchSourceMultiple)
				logger.Info(nil, "Found multiple results for show {Show} ({Year})", show.Name, show.Year)
			} else if len(showSearchResults) == 1 {
				show.FetchSource = entity.FetchSourceTVMaze
				self.showRepo.UpdateFetchSource(ctx, show.ID, entity.FetchSourceTVMaze)
				logger.Info(nil, "Found perfect match for show {Show} ({Year}): {Match}", show.Name, show.Year, showSearchResults[0])
			} else {
				show.FetchSource = entity.FetchSourceEmpty
				self.showRepo.UpdateFetchSource(ctx, show.ID, entity.FetchSourceEmpty)
			}
		}
	}()
}

func (self *ShowMatcher) createShowMetadata(ctx context.Context, show *entity.Show, res provider.ShowSearchResult) (*entity.ShowMetadata, error) {
	mid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	description := ""
	if res.Show.Summary != nil {
		description = *res.Show.Summary
	}
	mediumImage := ""
	originalImage := ""
	if res.Show.Image != nil {
		mediumImage = res.Show.Image.Medium
		originalImage = res.Show.Image.Original
	}
	metadata := entity.ShowMetadata{
		ID:               mid,
		ShowID:           show.ID,
		Name:             res.Show.Name,
		Url:              res.Show.URL,
		Description:      description,
		MediumImageUrl:   mediumImage,
		OriginalImageUrl: originalImage,
		FetchSource:      entity.FetchSourceTVMaze,
		FetchID:          res.Show.ID,
	}
	err = self.showMetadataRepo.Create(ctx, &metadata)
	if err != nil {
		return nil, err
	}
	err = self.createShowSeasonsMetadata(ctx, show, &metadata)
	if err != nil {
		return nil, err
	}
	return &metadata, nil
}

func (self *ShowMatcher) createShowSeasonsMetadata(ctx context.Context, show *entity.Show, metadata *entity.ShowMetadata) error {
	seasonMetadatas, err := self.metadataProvider.SearchSeasons(metadata.FetchID)
	if err != nil {
		return err
	}

	for _, sm := range seasonMetadatas {
		mid, err := uuid.NewV7()
		if err != nil {
			return err
		}
		mediumImage := ""
		originalImage := ""
		if sm.Image != nil {
			mediumImage = sm.Image.Medium
			originalImage = sm.Image.Original
		}
		m := entity.SeasonMetadata{
			ID:               mid,
			ShowID:           show.ID,
			ShowMetadataID:   metadata.ID,
			Url:              sm.Url,
			Number:           sm.Number,
			Summary:          sm.Summary,
			PremiereDate:     sm.PremiereDate,
			MediumImageUrl:   mediumImage,
			OriginalImageUrl: originalImage,
			FetchSource:      entity.FetchSourceTVMaze,
			FetchID:          sm.ID,
		}
		err = self.seasonMetadataRepo.Create(ctx, &m)
		if err != nil {
			return err
		}
		err = self.createSeasonEpisodeMetadata(ctx, show, &m)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *ShowMatcher) createSeasonEpisodeMetadata(ctx context.Context, show *entity.Show, metadata *entity.SeasonMetadata) error {
	episodeMetas, err := self.metadataProvider.SearchEpisodes(metadata.FetchID)
	if err != nil {
		return err
	}

	for _, em := range episodeMetas {
		mid, err := uuid.NewV7()
		if err != nil {
			return err
		}
		var _ = mid
		var _ = em
	}
	return nil
}
