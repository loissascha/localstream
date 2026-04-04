package background

import (
	"fmt"

	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/provider"
	"github.com/loissascha/localstream/internal/repository"
)

type ShowMatcher struct {
	Channel          chan *entity.Show
	metadataProvider provider.TVMetadataProvider
	showRepo         repository.ShowRepository
}

func NewShowMatcher(metadataProvider provider.TVMetadataProvider, showRepo repository.ShowRepository) *ShowMatcher {
	ch := make(chan *entity.Show)
	return &ShowMatcher{
		Channel:          ch,
		metadataProvider: metadataProvider,
		showRepo:         showRepo,
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

			if len(showSearchResults) > 1 {
				show.FetchSource = entity.FetchSourceMultiple
				// l.showRepo.UpdateFetchSource(context.Background(), show.ID, entity.FetchSourceMultiple)
				logger.Info(nil, "Found multiple results for show {Show} ({Year})", show.Name, show.Year)
			} else {
				show.FetchSource = entity.FetchSourceTVMaze
				// l.showRepo.UpdateFetchSource(context.Background(), show.ID, entity.FetchSourceTVMaze)
				logger.Info(nil, "Found perfect match for show {Show} ({Year}): {Match}", show.Name, show.Year, showSearchResults[0])
			}

			for _, res := range showSearchResults {
				fmt.Println(res)
			}
		}
	}()
}
