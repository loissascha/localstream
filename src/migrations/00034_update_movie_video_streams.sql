-- +goose Up
ALTER TABLE movie_video_streams ADD COLUMN container TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE movie_video_streams DROP COLUMN container;
