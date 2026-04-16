-- +goose Up
ALTER TABLE show_metadata
ADD COLUMN fetch_id INT NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE show_metadata
DROP COLUMN fetch_id;
