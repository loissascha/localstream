package handler

import (
	"net/http"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/provider"
	"github.com/loissascha/localstream/internal/service"
)

type SubtitleHandler struct {
	s               *server.Server
	authMiddleware  *middleware.AuthMiddleware
	subtitleService *service.SubtitleService
}

func NewSubtitleHandler(
	s *server.Server,
	authM *middleware.AuthMiddleware,
	subtitleService *service.SubtitleService,
) *SubtitleHandler {
	return &SubtitleHandler{
		s:               s,
		authMiddleware:  authM,
		subtitleService: subtitleService,
	}
}

func (h *SubtitleHandler) RegisterRoutes() {
	h.s.POSTI("/api/v1/subtitles/languages",
		h.languages,
		server.WithExportType[provider.SubtitleSupportedLanguage](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

func (h *SubtitleHandler) languages(w http.ResponseWriter, r *http.Request) {
	result, err := h.subtitleService.SupportedLanguages(r.Context())
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	respond.JSON(w, http.StatusOK, result)
}
