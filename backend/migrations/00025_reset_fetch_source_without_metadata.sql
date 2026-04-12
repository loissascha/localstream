-- +goose Up
UPDATE seasons s
SET fetch_source = 'none'
WHERE NOT EXISTS (
    SELECT 1
    FROM season_metadata sm
    WHERE sm.season_id = s.id
);

UPDATE episodes e
SET fetch_source = 'none'
WHERE NOT EXISTS (
    SELECT 1
    FROM episode_metadata em
    WHERE em.episode_id = e.id
);

-- +goose Down
SELECT 1;
