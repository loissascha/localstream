-- +goose Up
DELETE FROM season_metadata;

DROP INDEX IF EXISTS idx_season_metadata_show_id;
DROP INDEX IF EXISTS idx_season_metadata_show_id_number;

ALTER TABLE season_metadata
DROP COLUMN IF EXISTS show_id,
DROP COLUMN IF EXISTS show_metadata_id;

ALTER TABLE season_metadata
ADD COLUMN season_id UUID NOT NULL REFERENCES seasons(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_season_metadata_season_id ON season_metadata(season_id);

-- +goose Down
DROP INDEX IF EXISTS idx_season_metadata_season_id;

ALTER TABLE season_metadata
DROP COLUMN IF EXISTS season_id;

ALTER TABLE season_metadata
ADD COLUMN show_id UUID NOT NULL REFERENCES shows(id) ON DELETE CASCADE,
ADD COLUMN show_metadata_id UUID NOT NULL REFERENCES show_metadata(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_season_metadata_show_id ON season_metadata(show_id);
CREATE INDEX IF NOT EXISTS idx_season_metadata_show_id_number ON season_metadata(show_id, number);
