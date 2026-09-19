package background

import (
	"log/slog"
	"time"
)

type BackgroundService struct {
}

func NewBackgroundService() *BackgroundService {
	return &BackgroundService{}
}

func (s *BackgroundService) RunBackground() {
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
	// fetch all existing movies
	// fetch all existing shows
	// fetch all existing libraries
	return nil
}
