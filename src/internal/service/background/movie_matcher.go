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

type MovieMatcher struct {
	Channel          chan *entity.Movie
	metadataProvider provider.MovieMetadataProvider
	movieRepo        repository.MovieRepository
	movieMetaRepo    repository.MovieMetadataRepository
	movieMetaService *service.MovieMetadataService
}

func NewMovieMatcher(
	metadataProvider provider.MovieMetadataProvider,
	movieRepo repository.MovieRepository,
	movieMetaRepo repository.MovieMetadataRepository,
	movieMetaService *service.MovieMetadataService,
) *MovieMatcher {
	ch := make(chan *entity.Movie)
	return &MovieMatcher{
		Channel:          ch,
		metadataProvider: metadataProvider,
		movieRepo:        movieRepo,
		movieMetaRepo:    movieMetaRepo,
		movieMetaService: movieMetaService,
	}
}

func (self *MovieMatcher) RunBackground() {
	go func() {
		for {
			movie := <-self.Channel
			if !movie.FetchSource.IsNone() {
				continue
			}

			logger.Debug(nil, "____________ MOVIE ____________")
			logger.Debug(nil, "Name: {Name} | Year: {Year}", movie.Name, movie.Year)

			result, err := self.metadataProvider.SearchMovie(movie.Name, movie.Year)
			if err != nil {
				logger.Error(err, "Error getting movie data")
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			exactMatches := self.hasExactMatches(movie, result)
			// if exactly one exact match -> create only the metadata for that one
			if len(exactMatches) == 1 {
				self.movieMetaService.CreateMovieMetadata(ctx, movie, exactMatches[0])
				err := self.movieRepo.UpdateFetchSource(ctx, movie.ID, entity.FetchSourceTMDB)
				if err != nil {
					logger.Error(err, "Error Updating movie fetch source")
					continue
				}
				continue
			}

			for _, r := range result {
				self.movieMetaService.CreateMovieMetadata(ctx, movie, r)
			}

			// update movie fetch source based on amount of result
			if len(result) == 1 {
				err := self.movieRepo.UpdateFetchSource(ctx, movie.ID, entity.FetchSourceTMDB)
				if err != nil {
					logger.Error(err, "Error Updating movie fetch source")
					continue
				}
			} else if len(result) > 1 {
				err := self.movieRepo.UpdateFetchSource(ctx, movie.ID, entity.FetchSourceMultiple)
				if err != nil {
					logger.Error(err, "Error Updating movie fetch source")
					continue
				}
			} else {
				err := self.movieRepo.UpdateFetchSource(ctx, movie.ID, entity.FetchSourceEmpty)
				if err != nil {
					logger.Error(err, "Error Updating movie fetch source")
					continue
				}
			}
			logger.Debug(nil, "_______________________________")
		}
	}()
}

func (self *MovieMatcher) hasExactMatches(movie *entity.Movie, result []provider.MovieResult) []provider.MovieResult {
	exactMatches := []provider.MovieResult{}
	for _, r := range result {
		if r.Title == movie.Name {
			exactMatches = append(exactMatches, r)
		}
	}
	return exactMatches
}
