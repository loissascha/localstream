package handler

import (
	"net/http"

	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/middleware"
)

type ShowHandler struct {
	s              *server.Server
	authMiddleware *middleware.AuthMiddleware
}

func NewShowHandler(s *server.Server, authMiddleware *middleware.AuthMiddleware) *ShowHandler {
	return &ShowHandler{
		s:              s,
		authMiddleware: authMiddleware,
	}
}

func (h *ShowHandler) RegisterRoutes() {
	h.s.GET("/api/shows",
		h.listShows,
		// server.WithExportType[LibraryListItem](),
		// server.WithExportType[LibraryListResponse](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

func (h *ShowHandler) listShows(w http.ResponseWriter, r *http.Request) {

}
