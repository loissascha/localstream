package media

import (
	"log/slog"
	"time"

	"github.com/loissascha/localstream/internal/repository"
)

type IntegrityCheckService struct {
	libRepo repository.LibraryRepository

	showRepo    repository.ShowRepository
	seasonRepo  repository.SeasonRepository
	movieRepo   repository.MovieRepository
	episodeRepo repository.EpisodeRepository
}

func NewIntegrityCheckService(
	libRepo repository.LibraryRepository,

	showRepo repository.ShowRepository,
	seasonRepo repository.SeasonRepository,
	movieRepo repository.MovieRepository,
	episodeRepo repository.EpisodeRepository,
) *IntegrityCheckService {
	return &IntegrityCheckService{
		libRepo: libRepo,

		showRepo:    showRepo,
		seasonRepo:  seasonRepo,
		movieRepo:   movieRepo,
		episodeRepo: episodeRepo,
	}
}

func (s *IntegrityCheckService) StartBackground() {
	go func() {
		for {
			err := s.RunOnce()
			if err != nil {
				slog.Error("error appeared while running background service v2", "err", err)
			}
			time.Sleep(60 * time.Second)
		}
	}()
}

func (s *IntegrityCheckService) RunOnce() error {

	// search for an episode or movie that hasn't had the integrity check yet

	// do the integrity check on that file

	// record findings

	return nil
}
