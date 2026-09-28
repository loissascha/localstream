-- +goose Up
ALTER TABLE movie_video_streams ADD COLUMN pixel_format TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE movie_video_streams DROP COLUMN pixel_format;
