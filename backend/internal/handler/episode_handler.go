package handler

import (
	"net/http"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/service"
)

type EpisodeHandler struct {
	s              *server.Server
	authMiddleware *middleware.AuthMiddleware
	episodeService *service.EpisodeService
}

func NewEpisodeHandler(s *server.Server, authMiddleware *middleware.AuthMiddleware, episodeService *service.EpisodeService) *EpisodeHandler {
	return &EpisodeHandler{
		s:              s,
		authMiddleware: authMiddleware,
		episodeService: episodeService,
	}
}

func (h *EpisodeHandler) RegisterRoutes() {
	h.s.GET("/api/episodes/{seasonID}",
		h.listEpisodes,
		server.WithExportType[EpisodeInfo](),
		server.WithExportType[EpisodeListResponse](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

type EpisodeInfo struct {
	ID     string `json:"id"`
	Number int    `json:"number"`
}

type EpisodeListResponse struct {
	Episodes []EpisodeInfo `json:"episodes"`
}

func (h *EpisodeHandler) listEpisodes(w http.ResponseWriter, r *http.Request) {
	seasonID := r.PathValue("seasonID")
	episodes, err := h.episodeService.ListBySeasonID(r.Context(), seasonID)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read episodes"})
		return
	}

	result := make([]EpisodeInfo, 0, len(episodes))
	for _, episode := range episodes {
		result = append(result, EpisodeInfo{
			ID:     episode.ID.String(),
			Number: episode.Number,
		})
	}

	respond.JSON(w, http.StatusOK, EpisodeListResponse{Episodes: result})
}
