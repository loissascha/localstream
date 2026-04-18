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

- `src/` - Go application root
- `src/cmd/server/` - main server entrypoint
- `src/internal/` - backend packages
- `src/migrations/` - database migrations
- `src/frontend/` - SvelteKit frontend
- `src/frontend/build/` - built frontend assets served by the Go server in production

## Quick start

Prerequisites:

- Go `1.25.3+`
- pnpm
- PostgreSQL

1) Configure backend env

```bash
cp src/.env.example src/.env
```

At minimum, set `DATABASE_URL` and `PORT` in `src/.env`.

2) Configure frontend env

```bash
cp src/frontend/.env.example src/frontend/.env
```

Set `VITE_API_URL` to your backend origin (default local dev: `http://localhost:42069`).

3) Install dependencies

```bash
# backend
cd src && go mod download

# frontend
cd frontend && pnpm install
```

4) Run backend

```bash
cd src
go run ./cmd/server
```

5) Run frontend

```bash
cd src/frontend
pnpm dev
```

Open the app at Vite's dev URL (typically `http://localhost:5173`).

## Useful commands

Backend (run from `src/`):

- `go run ./cmd/server`
- `go build ./...`
- `go test ./...`
- `go vet ./...`
- `gofmt -w .`

Frontend (run from `src/frontend/`):

- `pnpm dev`
- `pnpm check`
- `pnpm lint`
- `pnpm build`
- `pnpm format`

## Environment variables

Backend (`src/.env`):

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

Frontend (`src/frontend/.env`):

- `VITE_API_URL` (backend base URL used by frontend API calls)

## Production

In production, the frontend is built into `src/frontend/build` and served by the main Go server. There is no separate frontend server process.

1) Configure environment files

```bash
cp src/.env.example src/.env
```

Set at least `PORT` and `DATABASE_URL` in `src/.env`. Set `APP_ENV=production` and configure `ALLOWED_ORIGINS` as needed.

2) Install dependencies

```bash
cd src && go mod download
cd frontend && pnpm install
```

3) Build the frontend

```bash
cd src/frontend
pnpm build
```

4) Start the Go server

```bash
cd src
go run ./cmd/server
```

Or build a binary first:

```bash
cd src
go build -o localstream ./cmd/server
./localstream
```

The Go server serves the built frontend from `src/frontend/build`, serves the API from the same process, and falls back to `index.html` for client-side routes.

## Notes

- Database migrations run automatically on backend startup.
- Type exports for frontend are generated in development by the backend server.
- In development, the frontend uses Vite and proxies `/api` requests to `VITE_API_URL`.

## License

MIT. See `LICENSE`.
