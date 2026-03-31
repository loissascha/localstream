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

type MovieHandler struct {
	s              *server.Server
	authMiddleware *middleware.AuthMiddleware
	movieService   *service.MovieService
}

func NewMovieHandler(s *server.Server, authM *middleware.AuthMiddleware, movieS *service.MovieService) *MovieHandler {
	return &MovieHandler{
		s:              s,
		authMiddleware: authM,
		movieService:   movieS,
	}
}

func (h *MovieHandler) RegisterRoutes() {
	h.s.GETI("/api/v1/movies/list",
		h.list,
		server.WithExportType[MovieListResponse](),
		server.WithExportType[MovieInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.GET("/api/movies/stream",
		h.streamVideo,
		server.WithMiddlewares(h.authMiddleware.RequireAuthURLToken),
	)
}

func (h *MovieHandler) streamVideo(w http.ResponseWriter, r *http.Request) {
	movieID := strings.TrimSpace(r.URL.Query().Get("id"))
	if movieID == "" {
		http.Error(w, "missing id query parameter", http.StatusBadRequest)
		return
	}

	movie, err := h.movieService.GetById(r.Context(), movieID)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read movie: " + err.Error()})
		return
	}
	if movie == nil {
		respond.JSON(w, http.StatusNotFound, map[string]string{"error": "movie not found"})
		return
	}

	file, err := os.Open(movie.Path)
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

	ext := strings.ToLower(filepath.Ext(movie.Path))
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

func (h *MovieHandler) list(w http.ResponseWriter, r *http.Request) {
	movies, err := h.movieService.List(r.Context())
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read movies: " + err.Error()})
		return
	}

	ml := []MovieInfo{}
	for _, m := range movies {
		ml = append(ml, toMovieInfo(&m))
	}

	respond.JSON(w, http.StatusOK, MovieListResponse{Movies: ml})
}
