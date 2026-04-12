package entity

import "github.com/google/uuid"

type AppSettings struct {
	ID                            uuid.UUID `db:"id"`
	ExecuteLibraryWatcher         bool      `db:"execute_library_watcher"`
	LibraryWatcherIntervalSeconds int       `db:"library_watcher_interval_seconds"`
}
