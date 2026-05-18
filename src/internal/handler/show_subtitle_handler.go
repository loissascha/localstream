package handler

import (
	"net/http"
	"strings"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/provider"
	"github.com/loissascha/localstream/internal/service"
)

type ShowSubtitleHandler struct {
	s               *server.Server
	authMiddleware  *middleware.AuthMiddleware
	subtitleService *service.ShowSubtitleService
}

func NewShowSubtitleHandler(
	s *server.Server,
	authM *middleware.AuthMiddleware,
	subtitleService *service.ShowSubtitleService,
) *ShowSubtitleHandler {
	return &ShowSubtitleHandler{
		s:               s,
		authMiddleware:  authM,
		subtitleService: subtitleService,
	}
}

func (h *ShowSubtitleHandler) RegisterRoutes() {
	h.s.POSTI("/api/v1/show/subtitles/search",
		h.searchByTerm,
		server.WithExportType[provider.SubtitleProviderResult](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

func (h *ShowSubtitleHandler) searchByTerm(w http.ResponseWriter, r *http.Request) {
	show := strings.TrimSpace(r.URL.Query().Get("show"))
	season := strings.TrimSpace(r.URL.Query().Get("season"))
	episode := strings.TrimSpace(r.URL.Query().Get("episode"))
	lang := strings.TrimSpace(r.URL.Query().Get("lang"))

	result, err := h.subtitleService.SearchByEpisode(r.Context(), show, season, episode, lang)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respond.JSON(w, http.StatusOK, result)
}
