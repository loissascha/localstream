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

	shows, err := l.showRepo.All(ctx)
	if err != nil {
		return err
	}

	for _, show := range shows {
		err := l.RunForShow(ctx, &show)
		if err != nil {
			return err
		}
	}

	movies, err := l.movieRepo.All(ctx)
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
		err := l.showRepo.DeleteByID(ctx, show.ID)
		if err != nil {
			return err
		}
		return nil
	}

	seasons, err := l.seasonRepo.ListByShowID(ctx, show.ID)
	if err != nil {
		return err
	}

	for _, season := range seasons {
		err := l.RunForSeason(ctx, &season)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *LibraryUncataloguer) RunForSeason(ctx context.Context, season *entity.Season) error {
	path := season.Path
	if !isDir(path) {
		err := l.seasonRepo.DeleteByID(ctx, season.ID)
		if err != nil {
			return err
		}
		return nil
	}

	episodes, err := l.episodeRepo.ListBySeasonID(ctx, season.ID)
	if err != nil {
		return err
	}

	for _, episode := range episodes {
		err := l.RunForEpisode(ctx, &episode)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *LibraryUncataloguer) RunForEpisode(ctx context.Context, episode *entity.Episode) error {
	path := episode.Path
	if !isFile(path) {
		err := l.episodeRepo.DeleteByID(ctx, episode.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *LibraryUncataloguer) RunForMovie(ctx context.Context, movie *entity.Movie) error {
	path := movie.Path
	if !isFile(path) {
		err := l.movieRepo.DeleteByID(ctx, movie.ID)
		if err != nil {
			return err
		}
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
