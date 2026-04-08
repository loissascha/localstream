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

type MovieMatcher struct {
	Channel          chan *entity.Movie
	metadataProvider provider.MovieMetadataProvider
	movieRepo        repository.MovieRepository
	movieMetaRepo    repository.MovieMetadataRepository
}

func NewMovieMatcher(metadataProvider provider.MovieMetadataProvider, movieRepo repository.MovieRepository, movieMetaRepo repository.MovieMetadataRepository) *MovieMatcher {
	ch := make(chan *entity.Movie)
	return &MovieMatcher{
		Channel:          ch,
		metadataProvider: metadataProvider,
		movieRepo:        movieRepo,
		movieMetaRepo:    movieMetaRepo,
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
				self.createMovieMetadata(ctx, movie, exactMatches[0])
				err := self.movieRepo.UpdateFetchSource(ctx, movie.ID, entity.FetchSourceTMDB)
				if err != nil {
					logger.Error(err, "Error Updating movie fetch source")
					continue
				}
				continue
			}

			for _, r := range result {
				self.createMovieMetadata(ctx, movie, r)
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
				err := self.movieRepo.UpdateFetchSource(ctx, movie.ID, entity.FetchSourceNone)
				if err != nil {
					logger.Error(err, "Error Updating movie fetch source")
					continue
				}
			}
			logger.Debug(nil, "_______________________________")
		}
	}()
}

func (self *MovieMatcher) createMovieMetadata(ctx context.Context, movie *entity.Movie, r provider.MovieResult) error {
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

	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}

	err = self.movieMetaRepo.Create(ctx, &entity.MovieMetadata{
		ID:               uuid,
		MovieID:          movie.ID,
		Name:             r.OriginalTitle,
		Url:              "",
		Description:      r.Overview,
		MediumImageUrl:   posterLink,
		BackdropImageUrl: backdropLink,
		FetchSource:      entity.FetchSourceTMDB,
	})
	if err != nil {
		return err
	}

	logger.Debug(nil, "------------------------------- ")
	return nil
}

func (self *MovieMatcher) hasExactMatches(movie *entity.Movie, result []provider.MovieResult) []provider.MovieResult {
	exactMatches := []provider.MovieResult{}
	for _, r := range result {
		if r.OriginalTitle == movie.Name {
			exactMatches = append(exactMatches, r)
		}
	}
	return exactMatches
}
