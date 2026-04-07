package background

import (
	"fmt"

	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/provider"
)

type MovieMatcher struct {
	Channel          chan *entity.Movie
	metadataProvider provider.MovieMetadataProvider
}

func NewMovieMatcher(metadataProvider provider.MovieMetadataProvider) *MovieMatcher {
	ch := make(chan *entity.Movie)
	return &MovieMatcher{
		Channel:          ch,
		metadataProvider: metadataProvider,
	}
}

func (self *MovieMatcher) RunBackground() {
	go func() {
		for {
			movie := <-self.Channel
			fmt.Println("movie", movie)
			result, err := self.metadataProvider.SearchMovie(movie.Name, movie.Year)
			if err != nil {
				logger.Error(err, "Error getting movie data")
				continue
			}

			logger.Info(nil, "Got the data!")
			fmt.Println(result)
		}
	}()
}
