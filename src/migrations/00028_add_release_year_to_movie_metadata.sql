-- +goose Up
ALTER TABLE movie_metadata
    ADD COLUMN IF NOT EXISTS release_year INTEGER NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE movie_metadata
    DROP COLUMN IF EXISTS release_year;
