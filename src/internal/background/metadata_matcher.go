package background

import (
	"log/slog"
	"time"
)

func (s *BackgroundService) RunMetadataMatchers() error {
	slog.Info("starting all metadata matchers")

	start := time.Now()
	err := s.runMoviesMatcher()
	if err != nil {
		return err
	}
	slog.Info("running movies matcher complete", "duration", time.Since(start))

	start = time.Now()
	err = s.runShowsMatcher()
	if err != nil {
		return err
	}
	slog.Info("running shows matcher complete", "duration", time.Since(start))

	start = time.Now()
	err = s.runSeasonsMatcher()
	if err != nil {
		return err
	}
	slog.Info("running seasons matcher complete", "duration", time.Since(start))

	start = time.Now()
	err = s.runEpisodesMatcher()
	if err != nil {
		return err
	}
	slog.Info("running episodes matcher complete", "duration", time.Since(start))

	slog.Info("finished all metadata matchers without errors")
	return nil
}
