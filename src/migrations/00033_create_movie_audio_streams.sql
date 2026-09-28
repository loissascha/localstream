-- +goose Up
CREATE TABLE IF NOT EXISTS movie_audio_streams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    stream_index INT NOT NULL,
    profile TEXT NOT NULL DEFAULT '',
    codec_name TEXT NOT NULL DEFAULT '',
    codec_long_name TEXT NOT NULL DEFAULT '',
    duration DOUBLE PRECISION NOT NULL DEFAULT 0,
    tags JSONB NOT NULL DEFAULT '{}'::jsonb,
    bit_rate BIGINT NOT NULL DEFAULT 0,
    CONSTRAINT uq_movie_audio_streams_movie_stream_index UNIQUE (movie_id, stream_index)
);

CREATE INDEX IF NOT EXISTS idx_movie_audio_streams_movie_id ON movie_audio_streams(movie_id);

-- +goose Down
DROP TABLE IF EXISTS movie_audio_streams;
