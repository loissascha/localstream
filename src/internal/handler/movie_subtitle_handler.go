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

type MovieSubtitleHandler struct {
	s                    *server.Server
	authMiddleware       *middleware.AuthMiddleware
	movieSubtitleService *service.MovieSubtitleService
}

func NewMovieSubtitleHandler(
	s *server.Server,
	authM *middleware.AuthMiddleware,
	movieSubtitleService *service.MovieSubtitleService,
) *MovieSubtitleHandler {
	return &MovieSubtitleHandler{
		s:                    s,
		authMiddleware:       authM,
		movieSubtitleService: movieSubtitleService,
	}
}

func (h *MovieSubtitleHandler) RegisterRoutes() {
	h.s.GETI("/api/v1/movie/subtitles/{movieID}",
		h.listByMovie,
		server.WithExportType[MovieSubtitleInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.POSTI("/api/v1/movie/subtitles/search",
		h.searchByTerm,
		server.WithExportType[provider.SubtitleProviderResult](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

func (h *MovieSubtitleHandler) searchByTerm(w http.ResponseWriter, r *http.Request) {
	term := strings.TrimSpace(r.URL.Query().Get("q"))
	result, err := h.movieSubtitleService.SearchByTerm(r.Context(), term)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respond.JSON(w, http.StatusOK, result)
}

func (h *MovieSubtitleHandler) listByMovie(w http.ResponseWriter, r *http.Request) {
	movieId := r.PathValue("movieID")
	subtitles, err := h.movieSubtitleService.ListByMovieID(r.Context(), movieId)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	result := []MovieSubtitleInfo{}
	for _, ms := range subtitles {
		result = append(result, toMovieSubtitleInfo(&ms))
	}

	respond.JSON(w, http.StatusOK, result)
}
