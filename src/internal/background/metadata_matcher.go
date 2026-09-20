package background

func (s *BackgroundService) RunMetadataMatchers() error {
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
	return nil
}
