-- +goose Up
CREATE TABLE IF NOT EXISTS movie_metadata (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    name TEXT NOT NULL DEFAULT '',
    url TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    medium_image_url TEXT NOT NULL DEFAULT '',
    backdrop_image_url TEXT NOT NULL DEFAULT '',
    fetch_source fetch_source NOT NULL DEFAULT 'none'
);

CREATE INDEX IF NOT EXISTS idx_movie_metadata_movie_id ON movie_metadata(movie_id);

-- +goose Down
DROP TABLE IF EXISTS movie_metadata;
