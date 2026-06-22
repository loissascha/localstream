package handler

import (
	"encoding/json"
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
	h.s.GETI("/api/v1/show/subtitles/{episodeID}",
		h.listByEpisode,
		server.WithExportType[SubtitleInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.POSTI("/api/v1/show/subtitles/search",
		h.searchByTerm,
		server.WithExportType[provider.SubtitleProviderResult](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.POSTI("/api/v1/show/subtitles/{episodeID}/create",
		h.createShowSubtitle,
		server.WithExportType[provider.SubtitleProviderResult](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

func (h *ShowSubtitleHandler) listByEpisode(w http.ResponseWriter, r *http.Request) {
	episodeId := r.PathValue("episodeID")
	subtitles, err := h.subtitleService.ListByEpisodeID(r.Context(), episodeId)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	result := []SubtitleInfo{}
	for _, ms := range subtitles {
		result = append(result, toSubtitleInfoEpisode(&ms))
	}

	respond.JSON(w, http.StatusOK, result)
}

func (h *ShowSubtitleHandler) createShowSubtitle(w http.ResponseWriter, r *http.Request) {
	episodeId := r.PathValue("episodeID")

	defer r.Body.Close()
	var result provider.SubtitleProviderResult
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		respond.JSON(w, http.StatusBadRequest, map[string]string{"json parsing error: ": err.Error()})
		return
	}

	if err := h.subtitleService.CreateFromSubtitleResult(r.Context(), episodeId, result); err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error: ": err.Error()})
		return
	}

	respond.JSON(w, http.StatusOK, true)
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
