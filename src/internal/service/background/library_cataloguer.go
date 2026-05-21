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
	movieMetaService   *service.MovieMetadataService
	showMetaService    *service.ShowMetadataService
	seasonMetaService  *service.SeasonMetadataService
	episodeMetaService *service.EpisodeMetadataService
	tvmetadataProvider provider.TVMetadataProvider
	showRepo           repository.ShowRepository
	seasonRepo         repository.SeasonRepository
	episodeRepo        repository.EpisodeRepository
	movieRepo          repository.MovieRepository
	showMatcher        *ShowMatcher
	movieMatcher       *MovieMatcher
	seasonMatcher      *SeasonMatcher
	episodeMatcher     *EpisodeMatcher
	allShows           []entity.Show
	allMovies          []entity.Movie
}

func NewLibraryCataloguer(
	libService *service.LibraryService,
	movieMetaService *service.MovieMetadataService,
	showRepo repository.ShowRepository,
	seasonRepo repository.SeasonRepository,
	episodeRepo repository.EpisodeRepository,
	movieRepo repository.MovieRepository,
	metadataProvider provider.TVMetadataProvider,
	movieMetadataProvider provider.MovieMetadataProvider,
	showMetadataRepo repository.ShowMetadataRepository,
	movieMetadataRepo repository.MovieMetadataRepository,
	seasonMetaRepo repository.SeasonMetadataRepository,
	episodeMetaRepo repository.EpisodeMetadataRepository,
	showMetaService *service.ShowMetadataService,
	seasonMetaService *service.SeasonMetadataService,
	episodeMetaService *service.EpisodeMetadataService,
) *LibraryCataloguer {
	showMatcher := NewShowMatcher(metadataProvider, showRepo, showMetadataRepo, seasonMetaRepo, episodeMetaRepo, showMetaService)
	showMatcher.RunBackground()

	movieMatcher := NewMovieMatcher(movieMetadataProvider, movieRepo, movieMetadataRepo, movieMetaService)
	movieMatcher.RunBackground()

	seasonMatcher := NewSeasonMatcher(metadataProvider, seasonMetaRepo, seasonRepo, showRepo, showMetadataRepo, seasonMetaService)
	seasonMatcher.RunBackground()

	episodeMatcher := NewEpisodeMatcher(metadataProvider, seasonMetaRepo, seasonRepo, showRepo, showMetadataRepo, episodeRepo, episodeMetaRepo, episodeMetaService)
	episodeMatcher.RunBackground()

	return &LibraryCataloguer{
		libService:         libService,
		tvmetadataProvider: tvmaze.NewTVMazeProvider(),
		showRepo:           showRepo,
		seasonRepo:         seasonRepo,
		episodeRepo:        episodeRepo,
		movieRepo:          movieRepo,
		showMatcher:        showMatcher,
		movieMatcher:       movieMatcher,
		seasonMatcher:      seasonMatcher,
		episodeMatcher:     episodeMatcher,
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var err error
	l.allShows, err = l.showRepo.All(ctx)
	if err != nil {
		return err
	}

	l.allMovies, err = l.movieRepo.All(ctx)
	if err != nil {
		return err
	}

	libraries, err := l.libService.List(ctx)
	if err != nil {
		return err
	}

	for _, lib := range libraries {
		err := l.RunLibrary(ctx, lib)
		if err != nil {
			return err
		}
	}

	l.allShows = nil
	l.allMovies = nil

	return nil
}

func (l *LibraryCataloguer) extractShows(basePath string, input []fResult) map[string]map[string][]string {
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

func (l *LibraryCataloguer) findOrCreateEpisode(ctx context.Context, seasonId uuid.UUID, episodeInfo *parsers.EpisodeInfo, basePath string, allExistingEpisodes []entity.Episode) (*entity.Episode, error) {
	p := path.Join(basePath, episodeInfo.RawName)

	var episode *entity.Episode
	for _, e := range allExistingEpisodes {
		if e.Path == p {
			episode = &e
			break
		}
	}

	if episode != nil {
		if episode.FetchSource.IsNone() {
			l.episodeMatcher.Channel <- episode
		}
		return episode, nil
	}

	episode = &entity.Episode{
		ID:          uuid.New(),
		SeasonID:    seasonId,
		Number:      episodeInfo.Episode,
		Path:        p,
		FetchSource: entity.FetchSourceNone,
	}

	err := l.episodeRepo.Create(ctx, episode)
	if err != nil {
		logger.Error(err, "Error creating episode")
		return nil, err
	}
	l.episodeMatcher.Channel <- episode
	return episode, nil
}

func (l *LibraryCataloguer) findOrCreateSeason(ctx context.Context, showId uuid.UUID, seasonInfo *parsers.SeasonInfo, basePath string, allExistingSeasons []entity.Season) (*entity.Season, error) {
	p := path.Join(basePath, seasonInfo.RawName)

	var season *entity.Season
	for _, s := range allExistingSeasons {
		if s.Path == p {
			season = &s
			break
		}
	}

	if season != nil {
		if season.FetchSource.IsNone() {
			l.seasonMatcher.Channel <- season
		}
		return season, nil
	}

	season = &entity.Season{
		ID:          uuid.New(),
		ShowID:      showId,
		Number:      seasonInfo.Season,
		Path:        p,
		FetchSource: entity.FetchSourceNone,
	}

	err := l.seasonRepo.Create(ctx, season)
	if err != nil {
		logger.Error(err, "Error creating season")
		return nil, err
	}
	l.seasonMatcher.Channel <- season
	return season, nil
}

func (l *LibraryCataloguer) findOrCreateShow(ctx context.Context, showInfo *parsers.ShowInfo, basePath string) (*entity.Show, error) {
	p := path.Join(basePath, showInfo.RawName)

	var show *entity.Show
	for _, sp := range l.allShows {
		if sp.Path == p {
			show = &sp
			break
		}
	}

	if show != nil {
		if show.FetchSource.IsNone() {
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

	err := l.showRepo.Create(ctx, show)
	if err != nil {
		logger.Error(err, "Error creating show")
		return nil, err
	}

	l.showMatcher.Channel <- show
	return show, nil
}

func (l *LibraryCataloguer) RunLibrary(ctx context.Context, library entity.Library) error {
	results, err := getAllFilesWithPath(library.Path, []string{"mp4"}) // "mkv" ? 
	if err != nil {
		return err
	}
	switch library.LibraryType {
	case entity.LibraryTypeShows:
		err := l.RunShowsLibrary(ctx, &library, results)
		if err != nil {
			return err
		}
		break
	case entity.LibraryTypeMovies:
		err := l.RunMoviesLibrary(ctx, &library, results)
		if err != nil {
			return err
		}
		break
	}

	return nil
}

func (l *LibraryCataloguer) RunShowsLibrary(ctx context.Context, library *entity.Library, results []fResult) error {
	shows := l.extractShows(library.Path, results)

	for show, seasons := range shows {
		showInfo, ok := parsers.ParseShowFromName(show)
		if !ok {
			logger.Error(nil, "Can't parse show name: {Show}. ParseShowFromName failed", show)
			continue
		}

		show, err := l.findOrCreateShow(ctx, showInfo, library.Path)
		if err != nil {
			return err
		}

		// get a list of all existing seasons for this show (so l.findOrCreateSeason can use the existing ones to check if it needs to create a new one)
		allExistingSeasons, err := l.seasonRepo.ListByShowID(ctx, show.ID)
		if err != nil {
			return err
		}

		for season, episodes := range seasons {
			seasonInfo, ok := parsers.ParseSeasonFromName(season)
			if !ok {
				logger.Error(nil, "Can't parse season name: {Season}. ParseSeasonFromName failed", season)
				continue
			}

			season, err := l.findOrCreateSeason(ctx, show.ID, seasonInfo, show.Path, allExistingSeasons)
			if err != nil {
				return err
			}

			// get a list of all existing episodes for this season (so l.findOrCreateEpisode can use the existing ones to check if it needs to create a new one)
			allExistingEpisodes, err := l.episodeRepo.ListBySeasonID(ctx, season.ID)
			if err != nil {
				return err
			}

			for _, episode := range episodes {
				episodeInfo, ok := parsers.ParseEpisodeFromFilename(episode)
				if !ok {
					logger.Error(nil, "Can't parse episode name: {Episode}. ParseEpisodeFromFilename failed", episode)
					continue
				}

				_, err := l.findOrCreateEpisode(ctx, season.ID, episodeInfo, season.Path, allExistingEpisodes)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (l *LibraryCataloguer) RunMoviesLibrary(ctx context.Context, library *entity.Library, results []fResult) error {
	for _, f := range results {
		var movie *entity.Movie
		for _, m := range l.allMovies {
			if m.Path == f.Path {
				movie = &m
				break
			}
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
		err := l.movieRepo.Create(ctx, movie)
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

func getAllFilesWithPath(startPoint string, extensions []string) ([]fResult, error) {
	result := []fResult{}
	err := filepath.WalkDir(startPoint, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// logger.Debug(nil, "DIR: {Dir}", path)
		} else {
			for _, extension := range extensions {
				if strings.HasSuffix(path, extension) {
					result = append(result, fResult{
						Path: path,
						Name: d.Name(),
					})
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
