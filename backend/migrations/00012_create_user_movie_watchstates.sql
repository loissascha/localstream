-- +goose Up
CREATE TABLE IF NOT EXISTS user_movie_watchstates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    movie_id UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    position DOUBLE PRECISION NOT NULL DEFAULT 0,
    duration DOUBLE PRECISION NOT NULL DEFAULT 0,
    finished BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_user_movie_watchstates_user_movie ON user_movie_watchstates(user_id, movie_id);
CREATE INDEX IF NOT EXISTS idx_user_movie_watchstates_user_updated_at ON user_movie_watchstates(user_id, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_movie_watchstates_movie_id ON user_movie_watchstates(movie_id);

-- +goose Down
DROP TABLE IF EXISTS user_movie_watchstates;
