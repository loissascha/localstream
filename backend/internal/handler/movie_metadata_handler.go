package handler

import (
	"net/http"
	"strings"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/provider"
	"github.com/loissascha/localstream/internal/service"
)

type MovieMetadataHandler struct {
	s                *server.Server
	authMiddleware   *middleware.AuthMiddleware
	movieMetaService *service.MovieMetadataService
}

func NewMovieMetadataHandler(s *server.Server, authM *middleware.AuthMiddleware, movieMetaService *service.MovieMetadataService) *MovieMetadataHandler {
	return &MovieMetadataHandler{
		s:                s,
		authMiddleware:   authM,
		movieMetaService: movieMetaService,
	}
}

func (h *MovieMetadataHandler) RegisterRoutes() {
	h.s.GETI("/api/v1/movie/metadata/{movieID}",
		h.listByMovie,
		server.WithExportType[MovieMetadataInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.POSTI("/api/v1/movie/metadata/{movieID}/set/primary/{id}",
		h.setPrimary,
		server.WithMiddlewares(h.authMiddleware.RequireAuthAdmin),
	)

	h.s.POSTI("/api/v1/movie/metadata/search",
		h.search,
		server.WithExportType[provider.MovieResult](),
		server.WithMiddlewares(h.authMiddleware.RequireAuthAdmin),
	)
}

func (h *MovieMetadataHandler) search(w http.ResponseWriter, r *http.Request) {
	searchQuery := strings.TrimSpace(r.URL.Query().Get("q"))
	if searchQuery == "" {
		http.Error(w, "missing search parameter", http.StatusBadRequest)
		return
	}
	result, err := h.movieMetaService.Search(r.Context(), searchQuery)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	respond.JSON(w, http.StatusOK, result)
}

func (h *MovieMetadataHandler) setPrimary(w http.ResponseWriter, r *http.Request) {
	movieId := r.PathValue("movieID")
	id := r.PathValue("id")

	err := h.movieMetaService.SetPrimaryForMovieID(r.Context(), movieId, id)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respond.JSON(w, http.StatusOK, id)
}

func (h *MovieMetadataHandler) listByMovie(w http.ResponseWriter, r *http.Request) {
	movieId := r.PathValue("movieID")
	metadata, err := h.movieMetaService.GetByMovieID(r.Context(), movieId)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	result := []MovieMetadataInfo{}
	for _, m := range metadata {
		result = append(result, toMovieMetadataInfo(&m))
	}

	respond.JSON(w, http.StatusOK, result)
}
