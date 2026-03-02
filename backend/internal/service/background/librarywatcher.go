package service

import (
	"time"

	"github.com/loissascha/go-logger/logger"
)

type LibraryWatcher struct {
}

func NewLibraryWatcher() *LibraryWatcher {
	return &LibraryWatcher{}
}

func (l *LibraryWatcher) RunBackground() {
	go func() {
		for {
			l.RunOnce()
			time.Sleep(5 * time.Second)
		}
	}()
}

func (l *LibraryWatcher) RunOnce() {
	logger.Info(nil, "Library Watcher running...")
}
