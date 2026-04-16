# AGENTS Guide

## Working Roots
- The real app lives under `src/`, not repo root.
- Run Go/backend commands from `src/`.
- Run Svelte/frontend commands from `src/frontend/`.
- There is no root-level task runner or workspace manifest.

## Backend
- Main entrypoint: `src/cmd/server/main.go`.
- Usual commands from `src/`: `go mod download`, `go run ./cmd/server`, `go build ./...`, `go test ./...`, `go vet ./...`, `gofmt -w .`.
- Focused parser tests live in `src/internal/parsers/`. Useful examples: `go test ./internal/parsers -run '^TestParseMovieFromFilename$'` and `go test ./internal/parsers -run 'TestParseSeasonFromName/specials' -count=1`.
- `godotenv.Load()` is called in the server, so backend env comes from `src/.env` when running from `src/`.
- Backend startup runs goose migrations automatically from `src/migrations/`.
- In development, the backend generates frontend API types at `src/frontend/src/lib/types/export_types.ts`; do not hand-edit that file.
- The backend also serves the built frontend from `src/frontend/build`, so production-like verification needs a frontend build first.

## Frontend
- Usual commands from `src/frontend/`: `bun install`, `bun run dev`, `bun run check`, `bun run lint`, `bun run build`, `bun run format`.
- There is no frontend test runner configured in `src/frontend/package.json`.
- `src/frontend/svelte.config.js` uses `@sveltejs/adapter-static` with `fallback: 'index.html'`.
- `src/frontend/src/routes/+layout.ts` sets `ssr = false`; treat the app as client-rendered.
- `VITE_API_URL` is used both by client code (`src/frontend/src/lib/consts.ts`) and the Vite `/api` proxy (`src/frontend/vite.config.ts`). Default local backend is `http://localhost:42069`.

## Env Files
- Backend template: `src/.env.example`.
- Frontend template: `src/frontend/.env.example`.
- Backend requires `PORT` and `DATABASE_URL`; frontend requires `VITE_API_URL`.

## Structure
- Backend packages are organized under `src/internal/` with the main flow split across `handler/`, `middleware/`, `service/`, `repository/`, `provider/`, and `database/`.
- Frontend routes are grouped under `src/frontend/src/routes/(auth)` and `src/frontend/src/routes/(protected)`; API helpers live under `src/frontend/src/lib/api/`.

## Verification
- Backend changes: run `go test ./...` and `go vet ./...` from `src/`.
- Frontend changes: run `bun run check` and `bun run lint` from `src/frontend/`.
