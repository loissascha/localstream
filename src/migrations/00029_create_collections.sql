-- +goose Up
CREATE TABLE IF NOT EXISTS collections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS collection_movies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    collection_id UUID NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    movie_id UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS collection_shows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    collection_id UUID NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    show_id UUID NOT NULL REFERENCES shows(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_collections_user_id ON collections(user_id);

CREATE UNIQUE INDEX IF NOT EXISTS uq_collection_movies_collection_movie
    ON collection_movies(collection_id, movie_id);
CREATE INDEX IF NOT EXISTS idx_collection_movies_collection_id ON collection_movies(collection_id);
CREATE INDEX IF NOT EXISTS idx_collection_movies_movie_id ON collection_movies(movie_id);

CREATE UNIQUE INDEX IF NOT EXISTS uq_collection_shows_collection_show
    ON collection_shows(collection_id, show_id);
CREATE INDEX IF NOT EXISTS idx_collection_shows_collection_id ON collection_shows(collection_id);
CREATE INDEX IF NOT EXISTS idx_collection_shows_show_id ON collection_shows(show_id);

-- +goose Down
DROP TABLE IF EXISTS collection_shows;
DROP TABLE IF EXISTS collection_movies;
DROP TABLE IF EXISTS collections;
