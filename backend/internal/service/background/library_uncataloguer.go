package background

import (
	"context"
	"os"
	"time"

	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type LibraryUncataloguer struct {
	showRepo    repository.ShowRepository
	seasonRepo  repository.SeasonRepository
	episodeRepo repository.EpisodeRepository
	movieRepo   repository.MovieRepository
}

func NewLibraryUncataloguer(showRepo repository.ShowRepository, seasonRepo repository.SeasonRepository, episodeRepo repository.EpisodeRepository, movieRepo repository.MovieRepository) *LibraryUncataloguer {
	return &LibraryUncataloguer{
		showRepo:    showRepo,
		seasonRepo:  seasonRepo,
		episodeRepo: episodeRepo,
		movieRepo:   movieRepo,
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
		err := l.RunForShow(ctx, &show)
		if err != nil {
			return err
		}
	}

	movies, err := l.movieRepo.List(ctx)
	if err != nil {
		return err
	}

	for _, movie := range movies {
		err := l.RunForMovie(ctx, &movie)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *LibraryUncataloguer) RunForShow(ctx context.Context, show *entity.Show) error {
	path := show.Path
	if !isDir(path) {
		// TODO: delete show from db
		return nil
	}

	seasons, err := l.seasonRepo.ListByShowID(ctx, show.ID)
	if err != nil {
		return err
	}

	for _, season := range seasons {
		l.RunForSeason(ctx, &season)
	}
	return nil
}

func (l *LibraryUncataloguer) RunForSeason(ctx context.Context, season *entity.Season) error {
	path := season.Path
	if !isDir(path) {
		// TODO: delete season from db
		return nil
	}

	// TODO: check all the episodes
	return nil
}

func (l *LibraryUncataloguer) RunForEpisode(ctx context.Context, episode *entity.Episode) error {
	path := episode.Path
	if !isFile(path) {
		// TODO: delete episode from db
	}
	return nil
}

func (l *LibraryUncataloguer) RunForMovie(ctx context.Context, movie *entity.Movie) error {
	path := movie.Path
	if !isFile(path) {
		// TODO: deelte movie from db
	}
	return nil
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
