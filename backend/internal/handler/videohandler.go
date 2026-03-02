package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
)

type VideoHandler struct {
	s                 *server.Server
	libraryDir        string
	allowedExtensions map[string]bool
	authMiddleware    *AuthMiddleware
}

type VideoListItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	MimeType string `json:"mimeType"`
}

type VideoListResponse struct {
	Videos []VideoListItem `json:"videos"`
}

func NewVideoHandler(s *server.Server, authMiddleware *AuthMiddleware) *VideoHandler {
	libraryDir := os.Getenv("VIDEO_LIBRARY_DIR")
	if strings.TrimSpace(libraryDir) == "" {
		libraryDir = "./videos"
	}

	return &VideoHandler{
		s:                 s,
		libraryDir:        libraryDir,
		allowedExtensions: parseAllowedExtensions(os.Getenv("VIDEO_ALLOWED_EXTENSIONS")),
		authMiddleware:    authMiddleware,
	}
}

func (h *VideoHandler) RegisterHandlers() {
	h.s.GET("/api/videos",
		h.listVideos,
		server.WithExportType[VideoListResponse](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
	h.s.GET("/api/videos/stream",
		h.streamVideo,
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

func (h *VideoHandler) listVideos(w http.ResponseWriter, r *http.Request) {
	entries, err := os.ReadDir(h.libraryDir)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read video library"})
		return
	}

	videos := []VideoListItem{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if !h.allowedExtensions[ext] {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		mimeType := mime.TypeByExtension(ext)
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		videos = append(videos, VideoListItem{
			ID:       encodeVideoID(entry.Name()),
			Name:     entry.Name(),
			Size:     info.Size(),
			MimeType: mimeType,
		})
	}

	sort.Slice(videos, func(i, j int) bool {
		return strings.ToLower(videos[i].Name) < strings.ToLower(videos[j].Name)
	})

	respond.JSON(w, http.StatusOK, VideoListResponse{Videos: videos})
}

func (h *VideoHandler) streamVideo(w http.ResponseWriter, r *http.Request) {
	videoID := strings.TrimSpace(r.URL.Query().Get("id"))
	if videoID == "" {
		http.Error(w, "missing id query parameter", http.StatusBadRequest)
		return
	}

	relativePath, err := decodeVideoID(videoID)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	absoluteVideoPath, err := h.resolveVideoPath(relativePath)
	if err != nil {
		http.Error(w, "video not found", http.StatusNotFound)
		return
	}

	file, err := os.Open(absoluteVideoPath)
	if err != nil {
		http.Error(w, "video not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil || info.IsDir() {
		http.Error(w, "video not found", http.StatusNotFound)
		return
	}

	fileSize := info.Size()
	if fileSize <= 0 {
		http.Error(w, "empty file", http.StatusNotFound)
		return
	}

	ext := strings.ToLower(filepath.Ext(absoluteVideoPath))
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Type", contentType)

	rangeHeader := strings.TrimSpace(r.Header.Get("Range"))
	if rangeHeader == "" {
		w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
		w.WriteHeader(http.StatusOK)
		_, _ = io.Copy(w, file)
		return
	}

	start, end, err := parseSingleRange(rangeHeader, fileSize)
	if err != nil {
		w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", fileSize))
		http.Error(w, "requested range not satisfiable", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	chunkSize := end - start + 1
	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		http.Error(w, "failed to stream file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.Header().Set("Content-Length", strconv.FormatInt(chunkSize, 10))
	w.WriteHeader(http.StatusPartialContent)
	_, _ = io.CopyN(w, file, chunkSize)
}

func (h *VideoHandler) resolveVideoPath(relativePath string) (string, error) {
	cleanPath := filepath.Clean(relativePath)
	if cleanPath == "." || filepath.IsAbs(cleanPath) || strings.HasPrefix(cleanPath, "..") {
		return "", errors.New("invalid path")
	}

	ext := strings.ToLower(filepath.Ext(cleanPath))
	if !h.allowedExtensions[ext] {
		return "", errors.New("invalid extension")
	}

	absLibraryDir, err := filepath.Abs(h.libraryDir)
	if err != nil {
		return "", err
	}

	absVideoPath, err := filepath.Abs(filepath.Join(absLibraryDir, cleanPath))
	if err != nil {
		return "", err
	}

	if absVideoPath != absLibraryDir && !strings.HasPrefix(absVideoPath, absLibraryDir+string(os.PathSeparator)) {
		return "", errors.New("path traversal blocked")
	}

	return absVideoPath, nil
}

func parseAllowedExtensions(raw string) map[string]bool {
	allowed := map[string]bool{}
	if strings.TrimSpace(raw) == "" {
		allowed[".mp4"] = true
		return allowed
	}

	parts := strings.Split(raw, ",")
	for _, part := range parts {
		ext := strings.ToLower(strings.TrimSpace(part))
		if ext == "" {
			continue
		}
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		allowed[ext] = true
	}

	if len(allowed) == 0 {
		allowed[".mp4"] = true
	}

	return allowed
}

func parseSingleRange(rangeHeader string, fileSize int64) (int64, int64, error) {
	if fileSize <= 0 {
		return 0, 0, errors.New("invalid file size")
	}

	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return 0, 0, errors.New("invalid range unit")
	}

	rawRange := strings.TrimSpace(strings.TrimPrefix(rangeHeader, "bytes="))
	if rawRange == "" || strings.Contains(rawRange, ",") {
		return 0, 0, errors.New("multiple or empty ranges are not supported")
	}

	parts := strings.SplitN(rawRange, "-", 2)
	if len(parts) != 2 {
		return 0, 0, errors.New("malformed range")
	}

	startPart := strings.TrimSpace(parts[0])
	endPart := strings.TrimSpace(parts[1])

	if startPart == "" {
		suffixLength, err := strconv.ParseInt(endPart, 10, 64)
		if err != nil || suffixLength <= 0 {
			return 0, 0, errors.New("invalid suffix range")
		}
		if suffixLength > fileSize {
			suffixLength = fileSize
		}
		start := fileSize - suffixLength
		return start, fileSize - 1, nil
	}

	start, err := strconv.ParseInt(startPart, 10, 64)
	if err != nil || start < 0 || start >= fileSize {
		return 0, 0, errors.New("invalid start range")
	}

	if endPart == "" {
		return start, fileSize - 1, nil
	}

	end, err := strconv.ParseInt(endPart, 10, 64)
	if err != nil || end < start {
		return 0, 0, errors.New("invalid end range")
	}
	if end >= fileSize {
		end = fileSize - 1
	}

	return start, end, nil
}

func encodeVideoID(relativePath string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(relativePath))
}

func decodeVideoID(videoID string) (string, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(videoID)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
