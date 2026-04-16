-- +goose Up
CREATE TYPE fetch_source AS ENUM ('none', 'tmdb', 'tvmaze');

CREATE TABLE IF NOT EXISTS shows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    fetch_source fetch_source NOT NULL DEFAULT 'none'
);

CREATE TABLE IF NOT EXISTS seasons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    show_id UUID NOT NULL REFERENCES shows(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    fetch_source fetch_source NOT NULL DEFAULT 'none'
);

CREATE TABLE IF NOT EXISTS episodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    season_id UUID NOT NULL REFERENCES seasons(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    fetch_source fetch_source NOT NULL DEFAULT 'none'
);

CREATE INDEX IF NOT EXISTS idx_seasons_show_id ON seasons(show_id);
CREATE INDEX IF NOT EXISTS idx_episodes_season_id ON episodes(season_id);

-- +goose Down
DROP TABLE IF EXISTS episodes;
DROP TABLE IF EXISTS seasons;
DROP TABLE IF EXISTS shows;

DROP TYPE IF EXISTS fetch_source;
