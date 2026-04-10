# AGENTS Guide

## Repo Boundaries
- `backend/`: Go API server, PostgreSQL access, and goose migrations.
- `frontend/`: SvelteKit 5 app (SPA/static build).

## Working Directory Rules
- Run backend commands from `backend/`.
- Run frontend commands from `frontend/`.
- No root-level task runner is configured.

## Backend Commands (from `backend/`)
- Setup: `go mod download`
- Run server: `go run ./cmd/server`
- Build: `go build ./...`
- Verify: `go test ./...` and `go vet ./...`
- Format (when editing Go): `gofmt -w .`

## Focused Backend Tests
- Single parser test: `go test ./internal/parsers -run '^TestParseMovieFromFilename$'`
- Single subtest: `go test ./internal/parsers -run '^TestParseSeasonFromName$/specials$'`
- Disable test cache while iterating: `go test ./internal/parsers -run '^TestParseSeasonFromName$' -count=1`

## Frontend Commands (from `frontend/`)
- Install deps: `bun install`
- Dev server: `bun run dev`
- Type/Svelte checks: `bun run check`
- Lint (Prettier check + ESLint): `bun run lint`
- Build: `bun run build`
- No frontend test runner is configured in `frontend/package.json`.

## Env + Runtime Gotchas
- Copy `backend/.env.example` to `backend/.env` and `frontend/.env.example` to `frontend/.env`.
- Backend requires `PORT` and `DATABASE_URL`; migrations run automatically on server startup.
- Backend exports API types only in development (`APP_ENV=development`) to `frontend/src/lib/types/export_types.ts`.
- Do not hand-edit `frontend/src/lib/types/export_types.ts` unless changing generation flow.

## Frontend API Wiring
- App API clients use `VITE_API_URL` (`frontend/src/lib/consts.ts`) for backend calls.
- Vite dev proxy for relative `/api` requests is controlled by `VITE_BACKEND_ORIGIN` in `frontend/vite.config.ts` (default `http://localhost:42069`).
- `frontend/src/routes/+layout.ts` sets `ssr = false`; treat frontend as client-rendered.

## Agent Verification Expectations
- Backend changes: run `go test ./...` and `go vet ./...`.
- Frontend changes: run `bun run check` and `bun run lint`.
