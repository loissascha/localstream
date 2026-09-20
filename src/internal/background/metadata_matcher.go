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

func (s *BackgroundService) runShowsMatcher() error {
	// get all the shows from the repo that have no metadata (fetchsource none)
	return nil
}

func (s *BackgroundService) runSeasonsMatcher() error {
	// get all the seasons from the repo that have a show with a valid fetch source (not none/empty/multiple)
	return nil
}

func (s *BackgroundService) runEpisodesMatcher() error {
	// get all the episodes from the repo that have a season with a valid fetch source
	return nil
}
