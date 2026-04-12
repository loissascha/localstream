package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type AppSettingsRepository struct {
	db *sqlx.DB
}

func NewAppSettingsRepository(db *sqlx.DB) *AppSettingsRepository {
	return &AppSettingsRepository{db: db}
}

func (r *AppSettingsRepository) Get(ctx context.Context) (*entity.AppSettings, error) {
	const query = `
		SELECT id, execute_library_watcher, library_watcher_interval_seconds
		FROM app_settings
		LIMIT 1
	`

	var appSettings entity.AppSettings
	if err := r.db.GetContext(ctx, &appSettings, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrAppSettingsNotFound
		}
		return nil, fmt.Errorf("get app settings: %w", err)
	}

	return &appSettings, nil
}

func (r *AppSettingsRepository) Update(ctx context.Context, appSettings *entity.AppSettings) error {
	const query = `
		UPDATE app_settings
		SET execute_library_watcher = $1,
			library_watcher_interval_seconds = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, appSettings.ExecuteLibraryWatcher, appSettings.LibraryWatcherIntervalSeconds, appSettings.ID)
	if err != nil {
		return fmt.Errorf("update app settings: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update app settings rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrAppSettingsNotFound
	}

	return nil
}

var _ repository.AppSettingsRepository = (*AppSettingsRepository)(nil)
