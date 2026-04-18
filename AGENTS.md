# AGENTS Guide

## Working Roots
- The app root is `src/`, not repo root.
- Run Go commands from `src/`.
- Run frontend commands from `src/frontend/`.
- There is no root-level workspace/task runner.

## Backend
- Entrypoint: `src/cmd/server/main.go`.
- Common commands from `src/`: `go mod download`, `go run ./cmd/server`, `go build ./...`, `go test ./...`, `go vet ./...`, `gofmt -w .`.
- The server calls `godotenv.Load()`, so local backend env is read from `src/.env` when started from `src/`.
- Backend startup automatically runs goose migrations from `src/migrations/`.
- In development, the server generates frontend API types at `src/frontend/src/lib/types/export_types.ts`; do not hand-edit that file.
- Focused parser checks: `go test ./internal/parsers -run '^TestParseMovieFromFilename$'` and `go test ./internal/parsers -run 'TestParseSeasonFromName/specials' -count=1`.

## Frontend
- Use `pnpm`, not `bun` (`src/frontend/pnpm-lock.yaml` is the lockfile).
- Common commands from `src/frontend/`: `pnpm install`, `pnpm dev`, `pnpm check`, `pnpm lint`, `pnpm build`, `pnpm format`.
- `pnpm check` runs `svelte-kit sync` first.
- `pnpm lint` is `prettier --check . && eslint .`.
- There is no frontend test runner in `src/frontend/package.json`.
- `src/frontend/src/routes/+layout.ts` sets `ssr = false`; treat the app as client-rendered.
- `src/frontend/svelte.config.js` uses `@sveltejs/adapter-static` with `fallback: 'index.html'`.
- `VITE_API_URL` drives both client API calls and the Vite `/api` proxy; default local backend is `http://localhost:42069`.

## Env And Build Quirks
- Backend env template: `src/.env.example`. Frontend env template: `src/frontend/.env.example`.
- Backend requires `PORT` and `DATABASE_URL`. Frontend requires `VITE_API_URL`.
- The Go server serves the built frontend from `FRONTEND_APP_DIR`, which defaults to `./frontend/build` in `src/.env.example`.
- Production-like verification needs a frontend build first (`pnpm build` in `src/frontend/`) before starting the Go server.

## Structure
- Backend packages live under `src/internal/`; the main split is `handler/`, `middleware/`, `service/`, `repository/`, `provider/`, and `database/`.
- Frontend routes are grouped under `src/frontend/src/routes/(auth)` and `src/frontend/src/routes/(protected)`.
- Frontend API helpers live under `src/frontend/src/lib/api/`.

## Verification
- Backend changes: run `go test ./...` and `go vet ./...` from `src/`.
- Frontend changes: run `pnpm check` and `pnpm lint` from `src/frontend/`.
