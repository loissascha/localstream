package background

import (
	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/provider"
)

type ShowMatcher struct {
	Channel          chan *entity.Show
	metadataProvider provider.TVMetadataProvider
}

func NewShowMatcher(metadataProvider provider.TVMetadataProvider) *ShowMatcher {
	ch := make(chan *entity.Show)
	return &ShowMatcher{
		Channel:          ch,
		metadataProvider: metadataProvider,
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
				logger.Error(nil, "Found multiple results for show {Show} ({Year})", show.Name, show.Year)
				continue
			}

			logger.Info(nil, "Found perfect match for show {Show} ({Year}): {Match}", show.Name, show.Year, showSearchResults[0])
		}
	}()
}
