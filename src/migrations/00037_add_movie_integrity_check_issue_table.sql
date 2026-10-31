-- +goose Up
CREATE TABLE IF NOT EXISTS movie_integrity_check_issues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    issue_level INTEGER NOT NULL,
    issue TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_movie_integrity_check_issues_movie_id ON movie_integrity_check_issues(movie_id);

-- +goose Down
DROP INDEX idx_movie_integrity_check_issues_movie_id;
DROP TABLE IF EXISTS movie_integrity_check_issues;
