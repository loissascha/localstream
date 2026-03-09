package handler

import (
	"net/http"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/service"
)

type SeasonHandler struct {
	s              *server.Server
	authMiddleware *middleware.AuthMiddleware
	seasonService  *service.SeasonService
}

func NewSeasonHandler(s *server.Server, authMiddleware *middleware.AuthMiddleware, seasonService *service.SeasonService) *SeasonHandler {
	return &SeasonHandler{
		s:              s,
		authMiddleware: authMiddleware,
		seasonService:  seasonService,
	}
}

func (h *SeasonHandler) RegisterRoutes() {
	h.s.GET("/api/seasons/{showID}",
		h.listSeasons,
		server.WithExportType[SeasonInfo](),
		server.WithExportType[SeasonListResponse](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

type SeasonInfo struct {
	ID     string `json:"id"`
	Number int    `json:"number"`
}

type SeasonListResponse struct {
	Seasons []SeasonInfo `json:"seasons"`
}

func (h *SeasonHandler) listSeasons(w http.ResponseWriter, r *http.Request) {
	showID := r.PathValue("showID")
	seasons, err := h.seasonService.ListByShowID(r.Context(), showID)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read seasons"})
		return
	}

	result := make([]SeasonInfo, 0, len(seasons))
	for _, season := range seasons {
		result = append(result, SeasonInfo{
			ID:     encoders.EncodeUUID(season.ID),
			Number: season.Number,
		})
	}

	respond.JSON(w, http.StatusOK, SeasonListResponse{Seasons: result})
}
