-- +goose Up
ALTER TABLE shows 
DROP COLUMN description;

ALTER TYPE fetch_source
ADD VALUE 'multiple' AFTER 'none';

CREATE TABLE IF NOT EXISTS show_metadata (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    show_id UUID NOT NULL REFERENCES shows(id) ON DELETE CASCADE,
    url TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    medium_image_url TEXT NOT NULL DEFAULT '',
    original_image_url TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_show_metadata_show_id ON show_metadata(show_id);

-- +goose Down
DROP TABLE IF EXISTS show_metadata;

ALTER TABLE shows 
ADD COLUMN description TEXT NOT NULL DEFAULT '';
