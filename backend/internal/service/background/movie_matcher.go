package background

import (
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
			if !movie.FetchSource.IsNone() {
				continue
			}
			continue // TODO: temporary

			logger.Debug(nil, "____________ MOVIE ____________")
			logger.Debug(nil, "Name: {Name} | Year: {Year}", movie.Name, movie.Year)

			result, err := self.metadataProvider.SearchMovie(movie.Name, movie.Year)
			if err != nil {
				logger.Error(err, "Error getting movie data")
				continue
			}

			for _, r := range result {
				logger.Debug(nil, "------------ RESULT ------------ ")

				backdropLink := ""
				posterLink := ""
				if r.BackdropPath != "" {
					backdropLink = "https://image.tmdb.org/t/p/w780" + r.BackdropPath
				}
				if r.PosterPath != "" {
					posterLink = "https://image.tmdb.org/t/p/w500" + r.PosterPath
				}

				logger.Debug(nil, "Title: {Title}", r.OriginalTitle)
				logger.Debug(nil, "Description: {Desc}", r.Overview)
				logger.Debug(nil, "Backdrop Link: {URL}", backdropLink)
				logger.Debug(nil, "Poster Link: {URL}", posterLink)

				logger.Debug(nil, "------------------------------- ")
			}
			logger.Debug(nil, "_______________________________")
		}
	}()
}
