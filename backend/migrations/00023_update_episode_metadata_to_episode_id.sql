-- +goose Up
DELETE FROM episode_metadata;

DROP INDEX IF EXISTS idx_episode_metadata_show_id;
DROP INDEX IF EXISTS idx_episode_metadata_season_metadata_id;
DROP INDEX IF EXISTS idx_episode_metadata_show_id_number;

ALTER TABLE episode_metadata
DROP COLUMN IF EXISTS show_id,
DROP COLUMN IF EXISTS season_metadata_id;

ALTER TABLE episode_metadata
ADD COLUMN episode_id UUID NOT NULL REFERENCES episodes(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_episode_metadata_episode_id ON episode_metadata(episode_id);

-- +goose Down
DROP INDEX IF EXISTS idx_episode_metadata_episode_id;

ALTER TABLE episode_metadata
DROP COLUMN IF EXISTS episode_id;

ALTER TABLE episode_metadata
ADD COLUMN show_id UUID NOT NULL REFERENCES shows(id) ON DELETE CASCADE,
ADD COLUMN season_metadata_id UUID NOT NULL REFERENCES season_metadata(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_episode_metadata_show_id ON episode_metadata(show_id);
CREATE INDEX IF NOT EXISTS idx_episode_metadata_season_metadata_id ON episode_metadata(season_metadata_id);
CREATE INDEX IF NOT EXISTS idx_episode_metadata_show_id_number ON episode_metadata(show_id, number);
