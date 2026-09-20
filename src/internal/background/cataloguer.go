package background

import (
	"context"
	"log/slog"
	"path"
	"strings"
	"time"

	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/parsers"
)

func (s *BackgroundService) runCataloguers() error {
	slog.Info("starting all cataloguers")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// fetch all existing shows
	start := time.Now()
	allShows, err := s.showRepo.All(ctx)
	if err != nil {
		return err
	}
	slog.Info("fetching all existing shows complete", "duration", time.Since(start))

	// fetch all existing movies
	start = time.Now()
	allMovies, err := s.movieRepo.All(ctx)
	if err != nil {
		return err
	}
	slog.Info("fetching all existing movies complete", "duration", time.Since(start))

	// fetch all existing libraries
	start = time.Now()
	libraries, err := s.libRepo.List(ctx)
	if err != nil {
		return err
	}
	slog.Info("fetching all existing libraries complete", "duration", time.Since(start))

	// run for each library
	start = time.Now()
	for _, lib := range libraries {
		err := s.runFullLibraryCataloguer(ctx, lib, allMovies, allShows)
		if err != nil {
			return err
		}
	}
	slog.Info("running all libraries complete", "duration", time.Since(start))

	slog.Info("finished all cataloguers without errors")
	return nil
}

func (s *BackgroundService) runFullLibraryCataloguer(
	ctx context.Context,
	lib entity.Library,
	existingMovies []entity.Movie,
	existingShows []entity.Show,
) error {

	start := time.Now()
	results, err := getAllFilesWithExtensionInPath(lib.Path, []string{"mp4"}) // "mkv" ?
	if err != nil {
		return err
	}
	slog.Info("get all files in library complete.", "duration", time.Since(start), "library.Name", lib.Name, "library.Path", lib.Path)

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
		if exists, _ := s.findMovieWithPathInList(f.Path, existingMovies); exists {
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

		err = s.fetchMetadataForMovie(ctx, movie)
		if err != nil {
			slog.Warn("direct metadata fetch for movie failed... will run again in the movie matcher", "err", err)
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
		ok, showID := s.findShowWithPathInList(showPath, existingShows)
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
			ok, seasonID := s.findSeasonWithPathInList(seasonPath, existingSeasons)
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
				ok, _ = s.findEpisodeWithPathInList(episodePath, existingEpisodes)
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
