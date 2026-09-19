package background

import (
	"context"
	"log/slog"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/parsers"
)

func (s *BackgroundService) runCataloguers() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// fetch all existing shows
	allShows, err := s.showRepo.All(ctx)
	if err != nil {
		return err
	}

	// fetch all existing movies
	allMovies, err := s.movieRepo.All(ctx)
	if err != nil {
		return err
	}

	// fetch all existing libraries
	libraries, err := s.libRepo.List(ctx)
	if err != nil {
		return err
	}

	// run for each library
	for _, lib := range libraries {
		err := s.runLibraryCataloguer(ctx, lib, allMovies, allShows)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *BackgroundService) runLibraryCataloguer(
	ctx context.Context,
	lib entity.Library,
	existingMovies []entity.Movie,
	existingShows []entity.Show,
) error {
	results, err := getAllFilesWithPath(lib.Path, []string{"mp4"}) // "mkv" ?
	if err != nil {
		return err
	}
	switch lib.LibraryType {
	case entity.LibraryTypeShows:
		err := s.runShowsLibraryCataloguer(ctx, &lib, results, existingShows)
		if err != nil {
			return err
		}
		break
	case entity.LibraryTypeMovies:
		err := s.runMoviesLibraryCataloguer(ctx, &lib, results, existingMovies)
		if err != nil {
			return err
		}
		break
	}

	return nil
}

func (s *BackgroundService) runMoviesLibraryCataloguer(ctx context.Context, lib *entity.Library, results []fResult, existingMovies []entity.Movie) error {
	for _, f := range results {
		if exists, _ := s.movieWithPathExistsInList(f.Path, existingMovies); exists {
			continue
		}

		movieInfo, ok := parsers.ParseMovieFromFilename(f.Name)
		if !ok {
			slog.Error("Can't parse movie filename", "fName", f.Name)
			continue
		}

		year := 0
		if movieInfo.Year != nil {
			year = *movieInfo.Year
		}
		movie := &entity.Movie{
			Name:      movieInfo.Title,
			Year:      year,
			CreatedAt: time.Now().UTC(),
			Path:      f.Path,
		}

		err := s.movieRepo.Create(ctx, movie)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BackgroundService) runShowsLibraryCataloguer(ctx context.Context, lib *entity.Library, results []fResult, existingShows []entity.Show) error {
	showsList := s.extractShows(lib.Path, results)

	// go through all the shows
	for show, seasons := range showsList {
		// parse info
		showInfo, ok := parsers.ParseShowFromName(show)
		if !ok {
			logger.Error(nil, "Can't parse show name: {Show}. ParseShowFromName failed", show)
			continue
		}

		// create show if it doesn't exist yet
		showPath := path.Join(lib.Path, showInfo.RawName)
		var err error
		ok, showID := s.showWithPathExistsInList(showPath, existingShows)
		if !ok {
			showID, err = s.createShow(ctx, showInfo, showPath)
			if err != nil {
				return err
			}
		}

		// get all the existing season for that show
		existingSeasons, err := s.seasonRepo.ListByShowID(ctx, showID)
		if err != nil {
			return err
		}

		// go through all the seasons
		for season, episodes := range seasons {
			// parse info
			seasonInfo, ok := parsers.ParseSeasonFromName(season)
			if !ok {
				logger.Error(nil, "Can't parse season name: {Season}. ParseSeasonFromName failed", season)
				continue
			}

			// create season if it doesn't exist yet
			seasonPath := path.Join(showPath, seasonInfo.RawName)
			ok, seasonID := s.seasonWithPathExistsInList(seasonPath, existingSeasons)
			if !ok {
				seasonID, err = s.createSeason(ctx, seasonInfo, showID, seasonPath)
				if err != nil {
					return err
				}
			}

			// get all existing episodes for that season
			existingEpisodes, err := s.episodeRepo.ListBySeasonID(ctx, seasonID)
			if err != nil {
				return err
			}

			for _, episode := range episodes {
				episodeInfo, ok := parsers.ParseEpisodeFromFilename(episode)
				if !ok {
					logger.Error(nil, "Can't parse episode name: {Episode}. ParseEpisodeFromFilename failed", episode)
					continue
				}

				episodePath := path.Join(seasonPath, episodeInfo.RawName)
				ok, _ = s.episodeWithPathExistsInList(episodePath, existingEpisodes)
				if !ok {
					_, err = s.createEpisode(ctx, episodeInfo, seasonID, episodePath)
					if err != nil {
						return err
					}
				}
			}
		}

	}
	return nil
}

func (s *BackgroundService) createEpisode(ctx context.Context, episodeInfo *parsers.EpisodeInfo, seasonId uuid.UUID, episodePath string) (uuid.UUID, error) {

	episode := &entity.Episode{
		ID:          uuid.New(),
		SeasonID:    seasonId,
		Number:      episodeInfo.Episode,
		Path:        episodePath,
		FetchSource: entity.FetchSourceNone,
	}

	err := s.episodeRepo.Create(ctx, episode)
	if err != nil {
		logger.Error(err, "Error creating episode")
		return uuid.Nil, err
	}
	return episode.ID, nil
}

func (s *BackgroundService) createSeason(ctx context.Context, seasonInfo *parsers.SeasonInfo, showId uuid.UUID, seasonPath string) (uuid.UUID, error) {
	season := &entity.Season{
		ID:          uuid.New(),
		ShowID:      showId,
		Number:      seasonInfo.Season,
		Path:        seasonPath,
		FetchSource: entity.FetchSourceNone,
	}

	err := s.seasonRepo.Create(ctx, season)
	if err != nil {
		logger.Error(err, "Error creating season")
		return uuid.Nil, err
	}
	return season.ID, nil
}

func (s *BackgroundService) createShow(ctx context.Context, showInfo *parsers.ShowInfo, showPath string) (uuid.UUID, error) {
	showE := &entity.Show{
		ID:          uuid.New(),
		Name:        showInfo.Series,
		Year:        0,
		Path:        showPath,
		FetchSource: entity.FetchSourceNone,
	}

	if showInfo.Year != nil {
		showE.Year = *showInfo.Year
	}

	err := s.showRepo.Create(ctx, showE)
	if err != nil {
		logger.Error(err, "Error creating show")
		return uuid.Nil, err
	}
	return showE.ID, nil
}

func (s *BackgroundService) extractShows(basePath string, input []fResult) map[string]map[string][]string {
	res := map[string]map[string][]string{}

	for _, result := range input {
		pathStr := strings.TrimPrefix(result.Path, basePath)
		pathStr = strings.TrimPrefix(pathStr, "/")
		split := strings.SplitN(pathStr, "/", 3)
		if len(split) != 3 {
			logger.Error(nil, "For this path the format of [ShowName]/[SeasonName]/[Episode] is not met. Please correct it. Path: {Path}", pathStr)
			continue
		}

		r, ok := res[split[0]]
		if !ok {
			res[split[0]] = map[string][]string{split[1]: []string{split[2]}}
			continue
		}
		r[split[1]] = append(r[split[1]], split[2])
	}

	return res
}

func (s *BackgroundService) showWithPathExistsInList(path string, list []entity.Show) (bool, uuid.UUID) {
	for _, show := range list {
		if show.Path == path {
			return true, show.ID
		}
	}
	return false, uuid.Nil
}

func (s *BackgroundService) seasonWithPathExistsInList(path string, list []entity.Season) (bool, uuid.UUID) {
	for _, season := range list {
		if season.Path == path {
			return true, season.ID
		}
	}
	return false, uuid.Nil
}

func (s *BackgroundService) episodeWithPathExistsInList(path string, list []entity.Episode) (bool, uuid.UUID) {
	for _, e := range list {
		if e.Path == path {
			return true, e.ID
		}
	}
	return false, uuid.Nil
}

func (s *BackgroundService) movieWithPathExistsInList(path string, list []entity.Movie) (bool, uuid.UUID) {
	for _, m := range list {
		if m.Path == path {
			return true, m.ID
		}
	}
	return false, uuid.Nil
}
