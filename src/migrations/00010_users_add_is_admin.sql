-- +goose Up
ALTER TABLE users
ADD COLUMN is_admin boolean NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE libraries
DROP COLUMN IF EXISTS is_admin;
