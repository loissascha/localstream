package handler

import (
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/middleware"
)

type MovieHandler struct {
	s              *server.Server
	authMiddleware *middleware.AuthMiddleware
}

func NewMovieHandler(s *server.Server, authM *middleware.AuthMiddleware) *MovieHandler {
	return &MovieHandler{
		s:              s,
		authMiddleware: authM,
	}
}

func (h *MovieHandler) RegisterRoutes() {
}
