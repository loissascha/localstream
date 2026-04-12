-- +goose Up
CREATE TABLE IF NOT EXISTS app_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    execute_library_watcher BOOLEAN NOT NULL DEFAULT false,
    library_watcher_interval_seconds INTEGER NOT NULL DEFAULT 120,
    singleton BOOLEAN NOT NULL DEFAULT true CHECK (singleton),
    CONSTRAINT app_settings_singleton_unique UNIQUE (singleton),
    CONSTRAINT app_settings_interval_positive CHECK (library_watcher_interval_seconds > 0)
);

INSERT INTO app_settings (
    execute_library_watcher,
    library_watcher_interval_seconds
)
VALUES (
    true,
    60
)
ON CONFLICT (singleton) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS app_settings;
