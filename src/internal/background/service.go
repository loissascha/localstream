package background

import (
	"context"
	"log/slog"
	"time"

	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type BackgroundService struct {
	showRepo  repository.ShowRepository
	movieRepo repository.MovieRepository
	libRepo   repository.LibraryRepository
}

func NewBackgroundService(
	showRepo repository.ShowRepository,
	movieRepo repository.MovieRepository,
	libRepo repository.LibraryRepository,
) *BackgroundService {
	return &BackgroundService{
		libRepo:   libRepo,
		showRepo:  showRepo,
		movieRepo: movieRepo,
	}
}

func (s *BackgroundService) StartBackground() {
	for {
		err := s.RunOnce()
		if err != nil {
			slog.Error("error appeared while running background service v2", "err", err)
		}
		time.Sleep(60 * time.Second)
	}
}

func (s *BackgroundService) RunOnce() error {

	// run the cataloguers (library cataloguer)
	err := s.runCataloguers()
	if err != nil {
		return err
	}

	// run metadata matchers

	return nil
}

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
	return nil
}
