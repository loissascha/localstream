-- +goose Up
CREATE TABLE IF NOT EXISTS user_watchstates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    show_id UUID NOT NULL REFERENCES shows(id) ON DELETE CASCADE,
    season_id UUID NOT NULL REFERENCES seasons(id) ON DELETE CASCADE,
    episode_id UUID NOT NULL REFERENCES episodes(id) ON DELETE CASCADE,
    position DOUBLE PRECISION NOT NULL DEFAULT 0,
    duration DOUBLE PRECISION NOT NULL DEFAULT 0,
    finished BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_user_watchstates_user_episode ON user_watchstates(user_id, episode_id);
CREATE INDEX IF NOT EXISTS idx_user_watchstates_user_updated_at ON user_watchstates(user_id, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_watchstates_show_id ON user_watchstates(show_id);
CREATE INDEX IF NOT EXISTS idx_user_watchstates_season_id ON user_watchstates(season_id);
CREATE INDEX IF NOT EXISTS idx_user_watchstates_episode_id ON user_watchstates(episode_id);

-- +goose Down
DROP TABLE IF EXISTS user_watchstates;
