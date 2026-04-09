package background

import (
	"context"
	"time"

	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type LibraryUncataloguer struct {
	showRepo  repository.ShowRepository
	movieRepo repository.MovieRepository
}

func NewLibraryUncataloguer(showRepo repository.ShowRepository, movieRepo repository.MovieRepository) *LibraryUncataloguer {
	return &LibraryUncataloguer{
		showRepo:  showRepo,
		movieRepo: movieRepo,
	}
}

func (l *LibraryUncataloguer) RunBackground() {
	go func() {
		for {
			err := l.RunOnce()
			if err != nil {
				logger.Error(err, "Error running library unwatcher")
			}
			time.Sleep(160 * time.Second)
		}
	}()
}

func (l *LibraryUncataloguer) RunOnce() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	shows, err := l.showRepo.List(ctx)
	if err != nil {
		return err
	}

	for _, show := range shows {
		err := l.RunForShow(&show)
		if err != nil {
			return err
		}
	}

	movies, err := l.movieRepo.List(ctx)
	if err != nil {
		return err
	}

	for _, movie := range movies {
		err := l.RunForMovie(&movie)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *LibraryUncataloguer) RunForShow(show *entity.Show) error {
	return nil
}

func (l *LibraryUncataloguer) RunForMovie(movie *entity.Movie) error {
	return nil
}
