-- +goose Up
CREATE VIEW user_all_watchstates AS
SELECT
    id,
    user_id,
    'episode'::TEXT AS media_type,
    episode_id AS media_id,
    position,
    duration,
    finished,
    created_at,
    updated_at
FROM user_watchstates
UNION ALL
SELECT
    id,
    user_id,
    'movie'::TEXT AS media_type,
    movie_id AS media_id,
    position,
    duration,
    finished,
    created_at,
    updated_at
FROM user_movie_watchstates;

-- +goose Down
DROP VIEW IF EXISTS user_all_watchstates;
