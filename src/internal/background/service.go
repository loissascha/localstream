package background

import (
	"log/slog"
	"time"

	"github.com/loissascha/localstream/internal/repository"
)

type BackgroundService struct {
	showRepo    repository.ShowRepository
	seasonRepo  repository.SeasonRepository
	movieRepo   repository.MovieRepository
	episodeRepo repository.EpisodeRepository
	libRepo     repository.LibraryRepository
}

func NewBackgroundService(
	showRepo repository.ShowRepository,
	seasonRepo repository.SeasonRepository,
	movieRepo repository.MovieRepository,
	libRepo repository.LibraryRepository,
	episodeRepo repository.EpisodeRepository,
) *BackgroundService {
	return &BackgroundService{
		libRepo:     libRepo,
		showRepo:    showRepo,
		seasonRepo:  seasonRepo,
		movieRepo:   movieRepo,
		episodeRepo: episodeRepo,
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

	// run the "uncataloguers"

	// run metadata matchers
	err = s.RunMetadataMatchers()
	if err != nil {
		return err
	}

	// run the stuff that deletes files that are no longer in use (because the metadata entries do not exist anymore for example)

	return nil
}
