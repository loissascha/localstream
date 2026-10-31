-- +goose Up
ALTER TABLE movies ADD COLUMN integrity_checked BOOLEAN DEFAULT false;
CREATE INDEX idx_movies_integrity_checked ON movies (integrity_checked);

-- +goose Down
DROP INDEX idx_movies_integrity_checked;
ALTER TABLE movies DROP COLUMN integrity_checked;
