package handler

import (
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/middleware"
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
}
