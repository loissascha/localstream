package background

import (
	"context"
	"os"
	"time"

	"github.com/loissascha/localstream/internal/entity"
)

func (s *BackgroundService) runUncataloguers() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	shows, err := s.showRepo.All(ctx)
	if err != nil {
		return err
	}

	for _, show := range shows {
		err := s.runUncataloguerForShow(ctx, &show)
		if err != nil {
			return err
		}
	}

	movies, err := s.movieRepo.All(ctx)
	if err != nil {
		return err
	}

	for _, movie := range movies {
		err := s.runUncataloguerForMovie(ctx, &movie)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *BackgroundService) runUncataloguerForShow(ctx context.Context, show *entity.Show) error {
	path := show.Path
	if !isDir(path) {
		err := s.showRepo.DeleteByID(ctx, show.ID)
		if err != nil {
			return err
		}
		return nil
	}

	seasons, err := s.seasonRepo.ListByShowID(ctx, show.ID)
	if err != nil {
		return err
	}

	for _, season := range seasons {
		err := s.runUncataloguerForSeason(ctx, &season)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BackgroundService) runUncataloguerForSeason(ctx context.Context, season *entity.Season) error {
	path := season.Path
	if !isDir(path) {
		err := s.seasonRepo.DeleteByID(ctx, season.ID)
		if err != nil {
			return err
		}
		return nil
	}

	episodes, err := s.episodeRepo.ListBySeasonID(ctx, season.ID)
	if err != nil {
		return err
	}

	for _, episode := range episodes {
		err := s.runUncataloguerForEpisode(ctx, &episode)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BackgroundService) runUncataloguerForEpisode(ctx context.Context, episode *entity.Episode) error {
	path := episode.Path
	if !isFile(path) {
		err := s.episodeRepo.DeleteByID(ctx, episode.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BackgroundService) runUncataloguerForMovie(ctx context.Context, movie *entity.Movie) error {
	path := movie.Path
	if !isFile(path) {
		err := s.movieRepo.DeleteByID(ctx, movie.ID)
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
