package handler

import (
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

}
