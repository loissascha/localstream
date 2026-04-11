-- +goose Up
CREATE TABLE IF NOT EXISTS season_metadata (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    show_id UUID NOT NULL REFERENCES shows(id) ON DELETE CASCADE,
    url TEXT NOT NULL DEFAULT '',
    number INT NOT NULL DEFAULT 0,
    summary TEXT NOT NULL DEFAULT '',
    premiere_date TEXT NOT NULL DEFAULT '',
    medium_image_url TEXT NOT NULL DEFAULT '',
    original_image_url TEXT NOT NULL DEFAULT '',
    fetch_id INT NOT NULL DEFAULT 0,
    fetch_source fetch_source NOT NULL DEFAULT 'none'
);

CREATE INDEX IF NOT EXISTS idx_season_metadata_show_id ON season_metadata(show_id);
CREATE INDEX IF NOT EXISTS idx_season_metadata_show_id_number ON season_metadata(show_id, number);

-- +goose Down
DROP TABLE IF EXISTS season_metadata;
