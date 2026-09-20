package background

import "log/slog"

func (s *BackgroundService) RunMetadataMatchers() error {
	slog.Info("starting all metadata matchers")
	err := s.runMoviesMatcher()
	if err != nil {
		return err
	}

	err = s.runShowsMatcher()
	if err != nil {
		return err
	}

	err = s.runSeasonsMatcher()
	if err != nil {
		return err
	}

	err = s.runEpisodesMatcher()
	if err != nil {
		return err
	}

	slog.Info("finished all metadata matchers without errors")
	return nil
}
