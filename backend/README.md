# Backend

## Environment

Copy `.env.example` to `.env` and adjust if needed.

- `PORT`: API server port
- `DATABASE_URL`: PostgreSQL connection string (required)
- `DB_MAX_OPEN_CONNS`: max open DB connections (optional, default `25`)
- `DB_MAX_IDLE_CONNS`: max idle DB connections (optional, default `25`)
- `DB_CONN_MAX_LIFETIME`: connection lifetime in seconds (optional, default `300`)
- `DB_MIGRATIONS_DIR`: goose migrations directory (optional, default `./migrations`)
- `VIDEO_LIBRARY_DIR`: local directory that contains video files
- `VIDEO_ALLOWED_EXTENSIONS`: comma-separated allowed extensions (milestone default: `.mp4`)

Database migrations are run automatically on server startup using goose.

## Run

```bash
go run ./cmd/server
```
