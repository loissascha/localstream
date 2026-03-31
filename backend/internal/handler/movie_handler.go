package handler

import (
	"net/http"

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
