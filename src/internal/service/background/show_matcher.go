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

type ShowMatcher struct {
	Channel             chan *entity.Show
	metadataProvider    provider.TVMetadataProvider
	showRepo            repository.ShowRepository
	showMetadataRepo    repository.ShowMetadataRepository
	seasonMetadataRepo  repository.SeasonMetadataRepository
	episodeMetadataRepo repository.EpisodeMetadataRepository
	showMetadataService *service.ShowMetadataService
}

func NewShowMatcher(metadataProvider provider.TVMetadataProvider, showRepo repository.ShowRepository, showMetadataRepo repository.ShowMetadataRepository, seasonMetaRepo repository.SeasonMetadataRepository, episodeMetaRepo repository.EpisodeMetadataRepository, showMetaService *service.ShowMetadataService) *ShowMatcher {
	ch := make(chan *entity.Show)
	return &ShowMatcher{
		Channel:             ch,
		metadataProvider:    metadataProvider,
		showRepo:            showRepo,
		showMetadataRepo:    showMetadataRepo,
		seasonMetadataRepo:  seasonMetaRepo,
		episodeMetadataRepo: episodeMetaRepo,
		showMetadataService: showMetaService,
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

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			hasError := false
			for _, res := range showSearchResults {
				err := self.showMetadataService.CreateShowMetadata(ctx, show, &res.Show)
				if err != nil {
					logger.Error(err, "Error creating show metadata")
					hasError = true
				}
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
