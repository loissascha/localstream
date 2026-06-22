-- +goose Up
CREATE TABLE IF NOT EXISTS episode_subtitles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    episode_id UUID NOT NULL REFERENCES episodes(id) ON DELETE CASCADE,
    path TEXT NOT NULL,
    name TEXT NOT NULL,
    lang_short TEXT NOT NULL DEFAULT '',
    lang TEXT NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_episode_subtitles_path ON episode_subtitles(path);
CREATE INDEX IF NOT EXISTS idx_episode_subtitles_episode_id ON episode_subtitles(episode_id);

-- +goose Down
DROP TABLE IF EXISTS episode_subtitles;
