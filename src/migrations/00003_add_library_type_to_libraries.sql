-- +goose Up
CREATE TYPE library_type AS ENUM ('movies', 'shows');

ALTER TABLE libraries
ADD COLUMN library_type library_type NOT NULL DEFAULT 'movies';

-- +goose Down
ALTER TABLE libraries
DROP COLUMN IF EXISTS library_type;

DROP TYPE IF EXISTS library_type;
