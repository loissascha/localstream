# localstream

`localstream` is a personal/fun self-hosted streaming app for local video files.

It combines a Go backend and a SvelteKit frontend to scan local media libraries, stream videos with HTTP byte-range support, and track watch progress for movies and shows.

## Current scope

- Stream local video files without loading full files into memory
- Organize content by libraries, shows, seasons, episodes, and movies
- Track user watch progress and "continue watching"
- Fetch metadata from TVMaze (shows) and TMDB (movies)
- Basic JWT auth with user/admin areas

This is an actively evolving side project, not a production-hardened platform.

## Stack

- Backend: Go, PostgreSQL, goose migrations
- Frontend: SvelteKit 5

## Repository layout

- `backend/` - API server, migrations, media scanning/streaming logic
- `frontend/` - SvelteKit app

## Quick start

Prerequisites:

- Go `1.25.3+`
- bun
- PostgreSQL

1) Configure backend env

```bash
cp backend/.env.example backend/.env
```

At minimum, set `DATABASE_URL` and `PORT` in `backend/.env`.

2) Configure frontend env

```bash
cp frontend/.env.example frontend/.env
```

Set `VITE_API_URL` to your backend origin (default local dev: `http://localhost:42069`).

3) Install dependencies

```bash
# backend
cd backend && go mod download

# frontend
cd ../frontend && bun install
```

4) Run backend

```bash
cd backend
go run ./cmd/server
```

5) Run frontend

```bash
cd frontend
bun run dev
```

Open the app at Vite's dev URL (typically `http://localhost:5173`).

## Useful commands

Backend (run from `backend/`):

- `go run ./cmd/server`
- `go test ./...`
- `go vet ./...`

Frontend (run from `frontend/`):

- `bun run dev`
- `bun run check`
- `bun run lint`
- `bun run build`

## Environment variables

Backend (`backend/.env`):

- `PORT` (required)
- `DATABASE_URL` (required)
- `DB_MAX_OPEN_CONNS` (optional)
- `DB_MAX_IDLE_CONNS` (optional)
- `DB_CONN_MAX_LIFETIME` (optional)
- `VIDEO_LIBRARY_DIR` (optional, default `./videos`)
- `VIDEO_ALLOWED_EXTENSIONS` (optional, default `.mp4`)
- `APP_ENV` (optional, `development`/`production`)
- `TMDB_API_KEY` (optional; needed for TMDB movie metadata)
- `ALLOWED_ORIGINS` (required for production environment. set to '*' to allow all hosts or set to frontend url for proper CORS handling)

Frontend (`frontend/.env`):

- `VITE_API_URL` (backend base URL used by frontend API calls)

## Notes

- Database migrations run automatically on backend startup.
- Type exports for frontend are generated in development by the backend server.

## License

MIT. See `LICENSE`.
