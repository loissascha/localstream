-- +goose Up
CREATE TABLE IF NOT EXISTS movies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    path TEXT NOT NULL,
	year INT NOT NULL DEFAULT 0,
	description TEXT NOT NULL DEFAULT '',
    fetch_source fetch_source NOT NULL DEFAULT 'none',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- +goose Down
DROP TABLE IF EXISTS movies;
