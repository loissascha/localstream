-- +goose Up
CREATE TABLE IF NOT EXISTS movie_subtitles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    path TEXT NOT NULL,
    name TEXT NOT NULL,
    lang_short TEXT NOT NULL DEFAULT '',
    lang TEXT NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_movie_subtitles_path ON movie_subtitles(path);
CREATE INDEX IF NOT EXISTS idx_movie_subtitles_movie_id ON movie_subtitles(movie_id);

-- +goose Down
DROP TABLE IF EXISTS movie_subtitles;
