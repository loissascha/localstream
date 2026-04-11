-- +goose Up
ALTER TABLE season_metadata
ADD COLUMN show_metadata_id UUID NOT NULL REFERENCES show_metadata(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE season_metadata
DROP COLUMN show_metadata_id;
