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
	"github.com/loissascha/localstream/internal/service"
)

type LibraryWatcher struct {
	libService         *service.LibraryService
	tvmetadataProvider provider.TVMetadataProvider
}

func NewLibraryWatcher(libService *service.LibraryService) *LibraryWatcher {
	return &LibraryWatcher{
		libService:         libService,
		tvmetadataProvider: tvmaze.NewTVMazeProvider(),
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

func (l *LibraryWatcher) RunLibrary(library entity.Library) error {
	results, err := getAllFilesWithPath(library.Path, "mp4")
	if err != nil {
		return err
	}
	for _, result := range results {

		episodeInfo, ok := parsers.ParseEpisodeFromFilename(result.Name)
		if !ok {
			logger.Warning(nil, "Couldn't parse file name {Name}", result.Name)
			continue
		}

		switch library.LibraryType {
		case entity.LibraryTypeMovies:
			logger.Debug(nil, "Movie: {File}", episodeInfo)
		case entity.LibraryTypeShows:
			logger.Debug(nil, "Show: {File}", episodeInfo)
			searchResults, err := l.tvmetadataProvider.SearchSeries(episodeInfo)
			if err != nil {
				logger.Error(err, "Failed searching show metadata for {Series}", episodeInfo.Series)
				continue
			}

			if len(searchResults) == 0 {
				logger.Info(nil, "No metadata results for {Series}", episodeInfo.Series)
				continue
			}

			bestMatch := searchResults[0]
			logger.Debug(nil, "Best metadata match for {Series}: {Name} ({ID}) with score {Score}", episodeInfo.Series, bestMatch.Show.Name, bestMatch.Show.ID, bestMatch.Score)
		}
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
