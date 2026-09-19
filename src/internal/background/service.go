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
	return nil
}
