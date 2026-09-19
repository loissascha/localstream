package background

import (
	"context"
	"log/slog"
	"time"

	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/parsers"
)

func (s *BackgroundService) runCataloguers() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// fetch all existing shows
	allShows, err := s.showRepo.All(ctx)
	if err != nil {
		return err
	}

	// fetch all existing movies
	allMovies, err := s.movieRepo.All(ctx)
	if err != nil {
		return err
	}

	// fetch all existing libraries
	libraries, err := s.libRepo.List(ctx)
	if err != nil {
		return err
	}

	// run for each library
	for _, lib := range libraries {
		err := s.runLibraryCataloguer(ctx, lib, allMovies, allShows)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *BackgroundService) runLibraryCataloguer(
	ctx context.Context,
	lib entity.Library,
	existingMovies []entity.Movie,
	existingShows []entity.Show,
) error {
	results, err := getAllFilesWithPath(lib.Path, []string{"mp4"}) // "mkv" ?
	if err != nil {
		return err
	}
	switch lib.LibraryType {
	case entity.LibraryTypeShows:
		err := s.runShowsLibraryCataloguer(ctx, &lib, results, existingShows)
		if err != nil {
			return err
		}
		break
	case entity.LibraryTypeMovies:
		err := s.runMoviesLibraryCataloguer(ctx, &lib, results, existingMovies)
		if err != nil {
			return err
		}
		break
	}

	return nil
}

func (s *BackgroundService) runMoviesLibraryCataloguer(ctx context.Context, lib *entity.Library, results []fResult, existingMovies []entity.Movie) error {
	for _, f := range results {
		if s.movieWithPathExistsInList(f.Path, existingMovies) {
			continue
		}
		movieInfo, ok := parsers.ParseMovieFromFilename(f.Name)
		if !ok {
			slog.Error("Can't parse movie filename", "fName", f.Name)
			continue
		}
		year := 0
		if movieInfo.Year != nil {
			year = *movieInfo.Year
		}
		movie := &entity.Movie{
			Name:      movieInfo.Title,
			Year:      year,
			CreatedAt: time.Now().UTC(),
			Path:      f.Path,
		}
		err := s.movieRepo.Create(ctx, movie)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BackgroundService) runShowsLibraryCataloguer(ctx context.Context, lib *entity.Library, results []fResult, existingShows []entity.Show) error {
	return nil
}

func (s *BackgroundService) movieWithPathExistsInList(path string, list []entity.Movie) bool {
	for _, m := range list {
		if m.Path == path {
			return true
		}
	}
	return false
}
