-- +goose Up
DELETE FROM episode_metadata;
DELETE FROM season_metadata;
DELETE FROM show_metadata;

UPDATE episodes
SET fetch_source = 'none';

UPDATE seasons
SET fetch_source = 'none';

UPDATE shows
SET fetch_source = 'none';

-- +goose Down
SELECT 1;
