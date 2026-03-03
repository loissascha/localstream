package background

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/parsers"
	"github.com/loissascha/localstream/internal/provider"
	"github.com/loissascha/localstream/internal/provider/tvmaze"
	"github.com/loissascha/localstream/internal/repository"
	"github.com/loissascha/localstream/internal/service"
)

type LibraryWatcher struct {
	libService         *service.LibraryService
	tvmetadataProvider provider.TVMetadataProvider
	showRepo           repository.ShowRepository
}

func NewLibraryWatcher(libService *service.LibraryService, showRepo repository.ShowRepository) *LibraryWatcher {
	return &LibraryWatcher{
		libService:         libService,
		tvmetadataProvider: tvmaze.NewTVMazeProvider(),
		showRepo:           showRepo,
	}
}

func (l *LibraryWatcher) RunBackground() {
	go func() {
		for {
			err := l.RunOnce()
			if err != nil {
				logger.Error(err, "Error running library watcher")
			}
			time.Sleep(5 * time.Second)
		}
	}()
}

func (l *LibraryWatcher) RunOnce() error {
	logger.Info(nil, "Library Watcher running...")

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

func (l *LibraryWatcher) extractShows(basePath string, input []fResult) map[string]map[string][]string {
	res := map[string]map[string][]string{}

	for _, result := range input {
		pathStr := strings.TrimPrefix(result.Path, basePath)
		pathStr = strings.TrimPrefix(pathStr, "/")
		logger.Debug(nil, "PathSTr: {Path}", pathStr)
		split := strings.SplitN(pathStr, "/", 3)
		if len(split) != 3 {
			logger.Error(nil, "Wrong length!")
			continue
		}
		logger.Debug(nil, "Split: {S}", split)

		r, ok := res[split[0]]
		if !ok {
			res[split[0]] = map[string][]string{split[1]: []string{split[2]}}
			continue
		}
		r[split[1]] = append(r[split[1]], split[2])
	}

	return res
}

func (l *LibraryWatcher) findOrCreateShow(showInfo *parsers.ShowInfo) {

}

func (l *LibraryWatcher) RunLibrary(library entity.Library) error {
	results, err := getAllFilesWithPath(library.Path, "mp4")
	if err != nil {
		return err
	}

	if library.LibraryType == entity.LibraryTypeShows {
		shows := l.extractShows(library.Path, results)

		for show, seasons := range shows {

			showInfo, ok := parsers.ParseShowFromName(show)
			if !ok {
				logger.Error(nil, "Can't parse show name: {Show}. ParseShowFromName failed", show)
				continue
			}

			if showInfo.Year != nil {
				logger.Info(nil, "Show parsed. Name: {Name} | Year: {Year} | Amount Seasons: {Seasons}", showInfo.Series, *showInfo.Year, len(seasons))
			} else {
				logger.Info(nil, "Show parsed. Name: {Name} | Amount Seasons: {Season}", showInfo.Series, len(seasons))
			}

			l.findOrCreateShow(showInfo)

			for season, episodes := range seasons {
				seasonInfo, ok := parsers.ParseSeasonFromName(season)
				if !ok {
					logger.Error(nil, "Can't parse season name: {Season}. ParseSeasonFromName failed", season)
					continue
				}

				logger.Info(nil, "Season parsed. Number: {Number}", seasonInfo.Season)

				for _, episode := range episodes {
					episodeInfo, ok := parsers.ParseEpisodeFromFilename(episode)
					if !ok {
						logger.Error(nil, "Can't parse episode name: {Episode}. ParseEpisodeFromFilename failed", episode)
						continue
					}

					logger.Info(nil, "Episode parsed. Number: {Number}", episodeInfo.Episode)
				}
			}
		}

		// go through each show
		// check if the show (on that path) already exists
		// if not -> create it
		// if yes -> go on
	}

	// for _, result := range results {
	//
	// 	episodeInfo, ok := parsers.ParseEpisodeFromFilename(result.Name)
	// 	if !ok {
	// 		logger.Warning(nil, "Couldn't parse file name {Name}", result.Name)
	// 		continue
	// 	}
	//
	// 	switch library.LibraryType {
	// 	case entity.LibraryTypeMovies:
	// 		logger.Debug(nil, "Movie: {File}", episodeInfo)
	// 	case entity.LibraryTypeShows:
	// 		logger.Debug(nil, "Show: {File} | Path: {Path}", episodeInfo, result.Path)
	//
	// 		// searchResults, err := l.tvmetadataProvider.SearchSeries(episodeInfo)
	// 		// if err != nil {
	// 		// 	logger.Error(err, "Failed searching show metadata for {Series}", episodeInfo.Series)
	// 		// 	continue
	// 		// }
	// 		//
	// 		// if len(searchResults) == 0 {
	// 		// 	logger.Info(nil, "No metadata results for {Series}", episodeInfo.Series)
	// 		// 	continue
	// 		// }
	// 		//
	// 		// bestMatch := searchResults[0]
	// 		// logger.Debug(nil, "Best metadata match for {Series}: {Name} ({ID}) with score {Score}", episodeInfo.Series, bestMatch.Show.Name, bestMatch.Show.ID, bestMatch.Score)
	// 	}
	// }
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
