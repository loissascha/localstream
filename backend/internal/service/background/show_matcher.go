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

type ShowMatcher struct {
	Channel          chan *entity.Show
	metadataProvider provider.TVMetadataProvider
	showRepo         repository.ShowRepository
	showMetadataRepo repository.ShowMetadataRepository
}

func NewShowMatcher(metadataProvider provider.TVMetadataProvider, showRepo repository.ShowRepository, showMetadataRepo repository.ShowMetadataRepository) *ShowMatcher {
	ch := make(chan *entity.Show)
	return &ShowMatcher{
		Channel:          ch,
		metadataProvider: metadataProvider,
		showRepo:         showRepo,
		showMetadataRepo: showMetadataRepo,
	}
}

func (l *ShowMatcher) RunBackground() {
	go func() {
		for {
			show := <-l.Channel
			logger.Info(nil, "New ShowID triggered! {Show}", show)

			showSearchResults, err := l.metadataProvider.SearchShow(show.Name, show.Year)
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
			for _, res := range showSearchResults {
				fmt.Println(res)
				mid, err := uuid.NewV7()
				if err != nil {
					logger.Fatal(err, "Can't create uuid")
					hasError = true
					continue
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
				}
				err = l.showMetadataRepo.Create(ctx, &metadata)
				if err != nil {
					logger.Fatal(err, "Can't create show metadata")
					hasError = true
				}
			}

			if hasError {
				continue
			}

			if len(showSearchResults) > 1 {
				show.FetchSource = entity.FetchSourceMultiple
				l.showRepo.UpdateFetchSource(ctx, show.ID, entity.FetchSourceMultiple)
				logger.Info(nil, "Found multiple results for show {Show} ({Year})", show.Name, show.Year)
			} else {
				show.FetchSource = entity.FetchSourceTVMaze
				l.showRepo.UpdateFetchSource(ctx, show.ID, entity.FetchSourceTVMaze)
				logger.Info(nil, "Found perfect match for show {Show} ({Year}): {Match}", show.Name, show.Year, showSearchResults[0])
			}
		}
	}()
}
