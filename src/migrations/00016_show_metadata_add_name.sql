-- +goose Up
ALTER TABLE show_metadata
ADD COLUMN name TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE show_metadata
DROP COLUMN name;
