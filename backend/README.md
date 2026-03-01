# Backend

## Environment

Copy `.env.example` to `.env` and adjust if needed.

- `PORT`: API server port
- `VIDEO_LIBRARY_DIR`: local directory that contains video files
- `VIDEO_ALLOWED_EXTENSIONS`: comma-separated allowed extensions (milestone default: `.mp4`)

## Milestone 01 Endpoints

- `GET /api/videos`
- `GET /api/videos/stream?id=<video-id>`

## Run

```bash
go run ./cmd/server
```
