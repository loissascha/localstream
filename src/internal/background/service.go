package background

import (
	"log/slog"
	"time"

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
