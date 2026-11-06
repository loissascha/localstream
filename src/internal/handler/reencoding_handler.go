package handler

import (
	"net/http"

	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/media"
	"github.com/loissascha/localstream/internal/middleware"
)

type ReencodingHandler struct {
	s                 *server.Server
	authMiddleware    *middleware.AuthMiddleware
	reencodingService *media.ReencodingService
}

func NewReencodingHandler(
	s *server.Server,
	authM *middleware.AuthMiddleware,
	reencodingService *media.ReencodingService,
) *ReencodingHandler {
	return &ReencodingHandler{
		s:                 s,
		authMiddleware:    authM,
		reencodingService: reencodingService,
	}
}

func (h *ReencodingHandler) RegisterRoutes() {
	h.s.GETI("/api/v1/movies/reencodings",
		h.listMovies,
		// server.WithExportType[MovieListResponse](),
		// server.WithExportType[MovieInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

func (h *ReencodingHandler) listMovies(w http.ResponseWriter, r *http.Request) {

}
