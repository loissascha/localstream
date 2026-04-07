package background

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/parsers"
	"github.com/loissascha/localstream/internal/provider"
	"github.com/loissascha/localstream/internal/provider/tvmaze"
	"github.com/loissascha/localstream/internal/repository"
	"github.com/loissascha/localstream/internal/service"
)

type LibraryCataloguer struct {
	libService         *service.LibraryService
	tvmetadataProvider provider.TVMetadataProvider
	showRepo           repository.ShowRepository
	seasonRepo         repository.SeasonRepository
	episodeRepo        repository.EpisodeRepository
	movieRepo          repository.MovieRepository
	showMatcher        *ShowMatcher
	movieMatcher       *MovieMatcher
	showMetadataRepo   repository.ShowMetadataRepository
}

func NewLibraryCataloguer(libService *service.LibraryService, showRepo repository.ShowRepository, seasonRepo repository.SeasonRepository, episodeRepo repository.EpisodeRepository, movieRepo repository.MovieRepository, metadataProvider provider.TVMetadataProvider, showMetadataRepo repository.ShowMetadataRepository) *LibraryCataloguer {

	showMatcher := NewShowMatcher(metadataProvider, showRepo, showMetadataRepo)
	showMatcher.RunBackground()

	movieMatcher := NewMovieMatcher()
	movieMatcher.RunBackground()

	return &LibraryCataloguer{
		libService:         libService,
		tvmetadataProvider: tvmaze.NewTVMazeProvider(),
		showRepo:           showRepo,
		seasonRepo:         seasonRepo,
		episodeRepo:        episodeRepo,
		movieRepo:          movieRepo,
		showMatcher:        showMatcher,
		movieMatcher:       movieMatcher,
	}
}

func (l *LibraryCataloguer) RunBackground() {
	go func() {
		for {
			err := l.RunOnce()
			if err != nil {
				logger.Error(err, "Error running library watcher")
			}
			time.Sleep(120 * time.Second)
		}
	}()
}

func (l *LibraryCataloguer) RunOnce() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	libraries, err := l.libService.List(ctx)
	if err != nil {
		return err
	}

	for _, lib := range libraries {
		err := l.RunLibrary(lib)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *LibraryCataloguer) extractShows(basePath string, input []fResult) map[string]map[string][]string {
	res := map[string]map[string][]string{}

	for _, result := range input {
		pathStr := strings.TrimPrefix(result.Path, basePath)
		pathStr = strings.TrimPrefix(pathStr, "/")
		split := strings.SplitN(pathStr, "/", 3)
		if len(split) != 3 {
			logger.Error(nil, "Wrong length!")
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

func (l *LibraryCataloguer) findOrCreateEpisode(seasonId uuid.UUID, episodeInfo *parsers.EpisodeInfo, basePath string) (*entity.Episode, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	p := path.Join(basePath, episodeInfo.RawName)
	episode, err := l.episodeRepo.GetByPathAndSeasonID(ctx, p, seasonId)
	if err != nil {
		logger.Error(err, "Couldn't get episode")
		return nil, err
	}
	if episode != nil {
		return episode, nil
	}

	episode = &entity.Episode{
		ID:          uuid.New(),
		SeasonID:    seasonId,
		Number:      episodeInfo.Episode,
		Path:        p,
		FetchSource: entity.FetchSourceNone,
	}

	err = l.episodeRepo.Create(ctx, episode)
	if err != nil {
		logger.Error(err, "Error creating episode")
		return nil, err
	}
	return episode, nil
}

func (l *LibraryCataloguer) findOrCreateSeason(showId uuid.UUID, seasonInfo *parsers.SeasonInfo, basePath string) (*entity.Season, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	p := path.Join(basePath, seasonInfo.RawName)
	season, err := l.seasonRepo.GetByPathAndShowID(ctx, p, showId)
	if err != nil {
		logger.Error(err, "Couldn't get season")
		return nil, err
	}
	if season != nil {
		return season, nil
	}

	season = &entity.Season{
		ID:          uuid.New(),
		ShowID:      showId,
		Number:      seasonInfo.Season,
		Path:        p,
		FetchSource: entity.FetchSourceNone,
	}

	err = l.seasonRepo.Create(ctx, season)
	if err != nil {
		logger.Error(err, "Error creating season")
		return nil, err
	}
	return season, nil
}

func (l *LibraryCataloguer) findOrCreateShow(showInfo *parsers.ShowInfo, basePath string) (*entity.Show, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	p := path.Join(basePath, showInfo.RawName)
	show, err := l.showRepo.GetByPath(ctx, p)
	if err != nil {
		logger.Error(err, "Couldn't get show by path")
		return nil, err
	}
	if show != nil {
		if show.FetchSource == entity.FetchSourceNone {
			l.showMatcher.Channel <- show
		}
		return show, nil
	}

	show = &entity.Show{
		ID:          uuid.New(),
		Name:        showInfo.Series,
		Year:        0,
		Path:        p,
		FetchSource: entity.FetchSourceNone,
	}

	if showInfo.Year != nil {
		show.Year = *showInfo.Year
	}

	err = l.showRepo.Create(ctx, show)
	if err != nil {
		logger.Error(err, "Error creating show")
		return nil, err
	}

	l.showMatcher.Channel <- show
	return show, nil
}

func (l *LibraryCataloguer) RunLibrary(library entity.Library) error {
	results, err := getAllFilesWithPath(library.Path, "mp4")
	if err != nil {
		return err
	}
	switch library.LibraryType {
	case entity.LibraryTypeShows:
		err := l.RunShowsLibrary(&library, results)
		if err != nil {
			return err
		}
		break
	case entity.LibraryTypeMovies:
		err := l.RunMoviesLibrary(&library, results)
		if err != nil {
			return err
		}
		break
	}

	return nil
}

func (l *LibraryCataloguer) RunShowsLibrary(library *entity.Library, results []fResult) error {
	shows := l.extractShows(library.Path, results)

	for show, seasons := range shows {

		showInfo, ok := parsers.ParseShowFromName(show)
		if !ok {
			logger.Error(nil, "Can't parse show name: {Show}. ParseShowFromName failed", show)
			continue
		}

		show, err := l.findOrCreateShow(showInfo, library.Path)
		if err != nil {
			return err
		}

		// if showInfo.Year != nil {
		// 	logger.Info(nil, "Show parsed. Name: {Name} (ID: {ID}) | Year: {Year} | Amount Seasons: {Seasons}", showInfo.Series, show.ID.String(), *showInfo.Year, len(seasons))
		// } else {
		// 	logger.Info(nil, "Show parsed. Name: {Name} (ID: {ID}) | Amount Seasons: {Season}", showInfo.Series, show.ID.String(), len(seasons))
		// }

		for season, episodes := range seasons {
			seasonInfo, ok := parsers.ParseSeasonFromName(season)
			if !ok {
				logger.Error(nil, "Can't parse season name: {Season}. ParseSeasonFromName failed", season)
				continue
			}

			season, err := l.findOrCreateSeason(show.ID, seasonInfo, show.Path)
			if err != nil {
				return err
			}

			for _, episode := range episodes {
				episodeInfo, ok := parsers.ParseEpisodeFromFilename(episode)
				if !ok {
					logger.Error(nil, "Can't parse episode name: {Episode}. ParseEpisodeFromFilename failed", episode)
					continue
				}

				_, err := l.findOrCreateEpisode(season.ID, episodeInfo, season.Path)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (l *LibraryCataloguer) RunMoviesLibrary(library *entity.Library, results []fResult) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, f := range results {
		movie, err := l.movieRepo.GetByPath(ctx, f.Path)
		if err != nil {
			return err
		}
		if movie != nil {
			l.movieMatcher.Channel <- movie
			continue
		}
		movieInfo, ok := parsers.ParseMovieFromFilename(f.Name)
		if !ok {
			logger.Error(nil, "Can't parse movie filename")
			continue
		}
		year := 0
		if movieInfo.Year != nil {
			year = *movieInfo.Year
		}
		movie = &entity.Movie{
			Name:      movieInfo.Title,
			Year:      year,
			CreatedAt: time.Now().UTC(),
			Path:      f.Path,
		}
		err = l.movieRepo.Create(ctx, movie)
		if err != nil {
			return err
		}
		l.movieMatcher.Channel <- movie
	}
	return nil
}

type fResult struct {
	Name string
	Path string
}

func getAllFilesWithPath(startPoint string, extension string) ([]fResult, error) {
	result := []fResult{}
	err := filepath.WalkDir(startPoint, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// logger.Debug(nil, "DIR: {Dir}", path)
		} else {
			if strings.HasSuffix(path, extension) {
				result = append(result, fResult{
					Path: path,
					Name: d.Name(),
				})
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
