-- +goose Up
ALTER TYPE fetch_source
ADD VALUE 'empty' AFTER 'none';

-- +goose Down
