package handler

import (
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/middleware"
)

type ShowMetadataHandler struct {
	s              *server.Server
	authMiddleware *middleware.AuthMiddleware
}

func NewShowMetadataHandler(s *server.Server, authM *middleware.AuthMiddleware) *ShowMetadataHandler {
	return &ShowMetadataHandler{
		s:              s,
		authMiddleware: authM,
	}
}

func (h *ShowMetadataHandler) RegisterRoutes() {

}
