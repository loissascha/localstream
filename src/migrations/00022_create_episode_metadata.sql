-- +goose Up
CREATE TABLE IF NOT EXISTS episode_metadata (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    show_id UUID NOT NULL REFERENCES shows(id) ON DELETE CASCADE,
    season_metadata_id UUID NOT NULL REFERENCES season_metadata(id) ON DELETE CASCADE,
    url TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL DEFAULT '',
    number INT NOT NULL DEFAULT 0,
    summary TEXT NOT NULL DEFAULT '',
    medium_image_url TEXT NOT NULL DEFAULT '',
    original_image_url TEXT NOT NULL DEFAULT '',
    fetch_id INT NOT NULL DEFAULT 0,
    fetch_source fetch_source NOT NULL DEFAULT 'none'
);

CREATE INDEX IF NOT EXISTS idx_episode_metadata_show_id ON episode_metadata(show_id);
CREATE INDEX IF NOT EXISTS idx_episode_metadata_season_metadata_id ON episode_metadata(season_metadata_id);
CREATE INDEX IF NOT EXISTS idx_episode_metadata_show_id_number ON episode_metadata(show_id, number);

-- +goose Down
DROP TABLE IF EXISTS episode_metadata;
