package handler

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/service"
)

type EpisodeHandler struct {
	s              *server.Server
	authMiddleware *middleware.AuthMiddleware
	episodeService *service.EpisodeService
}

func NewEpisodeHandler(s *server.Server, authMiddleware *middleware.AuthMiddleware, episodeService *service.EpisodeService) *EpisodeHandler {
	return &EpisodeHandler{
		s:              s,
		authMiddleware: authMiddleware,
		episodeService: episodeService,
	}
}

func (h *EpisodeHandler) RegisterRoutes() {
	h.s.GET("/api/episodes/{seasonID}",
		h.listEpisodes,
		server.WithExportType[EpisodeInfo](),
		server.WithExportType[EpisodeListResponse](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.GET("/api/episodes/stream",
		h.streamVideo,
		server.WithMiddlewares(h.authMiddleware.RequireAuthURLToken),
	)
}

func (h *EpisodeHandler) streamVideo(w http.ResponseWriter, r *http.Request) {
	episodeID := strings.TrimSpace(r.URL.Query().Get("id"))
	if episodeID == "" {
		http.Error(w, "missing id query parameter", http.StatusBadRequest)
		return
	}

	episode, err := h.episodeService.GetByID(r.Context(), episodeID)
	if err != nil {
		http.Error(w, "episode not found", http.StatusNotFound)
		return
	}

	file, err := os.Open(episode.Path)
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

	ext := strings.ToLower(filepath.Ext(episode.Path))
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

func (h *EpisodeHandler) listEpisodes(w http.ResponseWriter, r *http.Request) {
	seasonID := r.PathValue("seasonID")
	episodes, err := h.episodeService.ListBySeasonID(r.Context(), seasonID)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read episodes"})
		return
	}

	result := make([]EpisodeInfo, 0, len(episodes))
	for _, episode := range episodes {
		result = append(result, toEpisodeInfo(&episode))
	}

	respond.JSON(w, http.StatusOK, EpisodeListResponse{Episodes: result})
}
