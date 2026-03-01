# Milestone 01: MP4 Browser Streaming

## Goal

Stream local MP4 files to a browser session without loading whole files into memory.

## Out of Scope

- Adaptive bitrate streaming (HLS/DASH)
- Transcoding
- Authentication and multi-user logic
- Metadata database
- Subtitles and thumbnails

## Architecture (Current)

- Backend: Go + `github.com/loissascha/go-http-server`
- Frontend: SvelteKit 5
- Storage: local filesystem (`VIDEO_LIBRARY_DIR`)
- Transport: HTTP byte ranges (`Range` header)

## API Contract

### `GET /api/videos`

Returns streamable local video files from the configured library directory.

Response shape:

```json
{
  "videos": [
    {
      "id": "base64url-id",
      "name": "movie.mp4",
      "size": 123456789,
      "mimeType": "video/mp4"
    }
  ]
}
```

### `GET /api/videos/stream?id=<video-id>`

Streams a video file using HTTP range requests.

Expected behavior:

- Supports `Range` requests and returns `206 Partial Content`
- Returns `200 OK` when no `Range` header is provided
- Returns `416 Requested Range Not Satisfiable` for invalid ranges
- Always sets `Accept-Ranges: bytes`

## Environment Variables

- `PORT` (required): backend port
- `VIDEO_LIBRARY_DIR` (optional): defaults to `./videos`
- `VIDEO_ALLOWED_EXTENSIONS` (optional): comma-separated list; defaults to `.mp4`

## Acceptance Criteria

- Large MP4 files start playback quickly
- Seeking to arbitrary timestamps works
- Backend does not read the full file into memory
- Invalid range requests are handled correctly with `416`
- Only allowed extensions are listed/streamed
- Path traversal is blocked

## Implementation Plan

1. Add backend video handler (`/api/videos`, `/api/videos/stream`)
2. Implement secure file resolution + range streaming
3. Add simple Svelte page with list + player
4. Configure frontend dev proxy for `/api` to backend
5. Verify with a local large MP4 file

## Manual Verification

1. Put one or more `.mp4` files into `backend/videos` (or set `VIDEO_LIBRARY_DIR`)
2. Start backend
3. Start frontend
4. Open app and select a file
5. Seek multiple times while checking network panel for `206` responses

## Next Milestone (Later)

- ffprobe metadata extraction
- thumbnail/poster generation
- watch progress persistence
- HLS transcoding pipeline
