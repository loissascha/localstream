package background

import (
	"context"
	"time"

	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/provider"
)

func (s *BackgroundService) runMoviesMatcher() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// get all the movies from the repo that have no metadata (fetchsource none)
	moviesWithoutMetadata, err := s.movieRepo.NecessaryForMetadataFetch(ctx)
	if err != nil {
		return err
	}

	for _, movie := range moviesWithoutMetadata {
		s.fetchMetadataForMovie(ctx, &movie)
	}
	return nil
}

func (s *BackgroundService) fetchMetadataForMovie(ctx context.Context, movie *entity.Movie) error {
	if !movie.FetchSource.IsNone() {
		return nil
	}

	logger.Debug(nil, "____________ MOVIE ____________")
	logger.Debug(nil, "Name: {Name} | Year: {Year}", movie.Name, movie.Year)

	result, err := s.movieMetadataProvider.SearchMovie(movie.Name, movie.Year)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// if exactly one exact match -> create only the metadata for that one
	exactMatches := s.movieMetadataHasExactMatches(movie, result)
	if len(exactMatches) == 1 {
		s.movieMetaService.CreateMovieMetadata(ctx, movie, exactMatches[0])
		err := s.movieRepo.UpdateFetchSource(ctx, movie.ID, entity.FetchSourceTMDB)
		if err != nil {
			return err
		}
		return nil
	}

	for _, r := range result {
		s.movieMetaService.CreateMovieMetadata(ctx, movie, r)
	}

	// update movie fetch source based on amount of result
	if len(result) == 1 {
		err := s.movieRepo.UpdateFetchSource(ctx, movie.ID, entity.FetchSourceTMDB)
		if err != nil {
			return err
		}
	} else if len(result) > 1 {
		err := s.movieRepo.UpdateFetchSource(ctx, movie.ID, entity.FetchSourceMultiple)
		if err != nil {
			return err
		}
	} else {
		err := s.movieRepo.UpdateFetchSource(ctx, movie.ID, entity.FetchSourceEmpty)
		if err != nil {
			return err
		}
	}
	logger.Debug(nil, "_______________________________")
	return nil
}

func (s *BackgroundService) movieMetadataHasExactMatches(movie *entity.Movie, result []provider.MovieResult) []provider.MovieResult {
	exactMatches := []provider.MovieResult{}
	for _, r := range result {
		if r.Title == movie.Name {
			exactMatches = append(exactMatches, r)
		}
	}
	return exactMatches
}
