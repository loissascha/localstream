-- +goose Up
ALTER TABLE show_metadata
ADD COLUMN fetch_source fetch_source NOT NULL DEFAULT 'none';

-- +goose Down
ALTER TABLE show_metadata
DROP COLUMN fetch_source;
