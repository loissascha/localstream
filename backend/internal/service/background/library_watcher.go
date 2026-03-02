package service

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/loissascha/go-logger/logger"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/service"
)

type LibraryWatcher struct {
	libService *service.LibraryService
}

func NewLibraryWatcher(libService *service.LibraryService) *LibraryWatcher {
	return &LibraryWatcher{
		libService: libService,
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

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)

	libraries, err := l.libService.List(ctx)
	if err != nil {
		return err
	}

	for _, lib := range libraries {
		err := l.runOnLibrary(lib)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *LibraryWatcher) runOnLibrary(library entity.Library) error {
	paths, err := getAllFilesWithPath(library.Path, "mp4")
	if err != nil {
		return err
	}
	for _, path := range paths {
		logger.Debug(nil, "{File}", path)
	}
	return nil
}

func getAllFilesWithPath(startPoint string, extension string) ([]string, error) {
	result := []string{}
	err := filepath.WalkDir(startPoint, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// logger.Debug(nil, "DIR: {Dir}", path)
		} else {
			if strings.HasSuffix(path, extension) {
				result = append(result, path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
