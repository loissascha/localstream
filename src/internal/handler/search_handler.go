package handler

import (
	"net/http"
	"strings"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/service"
)

type SearchHandler struct {
	s              *server.Server
	authMiddleware *middleware.AuthMiddleware
	showSerivce    *service.ShowService
	movieService   *service.MovieService
}

func NewSearchHandler(s *server.Server, authMiddleware *middleware.AuthMiddleware, showService *service.ShowService, movieService *service.MovieService) *SearchHandler {
	return &SearchHandler{
		s:              s,
		authMiddleware: authMiddleware,
		showSerivce:    showService,
		movieService:   movieService,
	}
}

func (h *SearchHandler) RegisterRoutes() {
	h.s.GETI("/api/v1/search",
		h.search,
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

func (h *SearchHandler) search(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if len(query) < 3 {
		respond.JSON(w, http.StatusNoContent, nil)
		return
	}
}
