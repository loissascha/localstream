package handler

import (
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/service"
)

type ShowMetadataHandler struct {
	s               *server.Server
	authMiddleware  *middleware.AuthMiddleware
	showMetaService *service.ShowMetadataService
}

func NewShowMetadataHandler(s *server.Server, authM *middleware.AuthMiddleware, showMetaService *service.ShowMetadataService) *ShowMetadataHandler {
	return &ShowMetadataHandler{
		s:               s,
		authMiddleware:  authM,
		showMetaService: showMetaService,
	}
}

func (h *ShowMetadataHandler) RegisterRoutes() {

}
