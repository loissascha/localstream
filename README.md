# localstream

`localstream` is a personal/fun self-hosted streaming app for local video files.

It combines a Go backend and a SvelteKit frontend to scan local media libraries, stream videos with HTTP byte-range support, and track watch progress for movies and shows.

## Current scope

- Stream local video files
- Track user watch progress
- Fetch metadata from TVMaze (shows) and TMDB (movies)

This is an actively evolving side project, not a production-hardened platform.

## Stack

- Backend: Go, PostgreSQL
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
This is only required for the development environment. In production the frontend is hosted by the backend and therefore the frontend doesn't need this.
If you're running the default PORT setup on the backend you don't need to create tne .env file.

3) Install dependencies

```bash
# backend
cd src && go mod download

# frontend
cd src/frontend && pnpm install
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

## Production

The Go server serves the built frontend from `src/frontend/build` and serves the API from the same process.

I chose this path because it's the easiest when the app is self hosted in a home network environment and someone tries to access it by for example Tailscale IP the frontend can still connect to the backend by its relative path.

## Notes

- Database migrations run automatically on backend startup.
- Type exports for frontend are generated in development by the backend server.
- In development, the frontend uses Vite and proxies `/api` requests to `VITE_API_URL`.

## Usage of AI 

This project was *not* vibe-coded or heavily influenced by AI. 

I used AI in a very controlled manner. 

For example: 
- using it as a replacement to "Googling something" 
- creating a set of svg icons in the frontend (most of them as placeholders that have already been replaced)
- writing documentation (parts of this README for example)
- letting it fill out some pre-defined boilerplate code (for example: creating the repository interface/struct for a new entity with basic CRUD operations)
- writing small (often temporary / later refactored) components (for example: the original VideoPlayer.svelte component; regex parsers in the backend to extract show/movie names, episode numbers, season numbers from raw path/file names; the original byte range video stream backend handler)

## License

MIT. See `LICENSE`.
