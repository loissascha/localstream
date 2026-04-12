package handler

import (
	"net/http"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
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
	h.s.GET("/api/v1/seasons/show/{showID}",
		h.listSeasons,
		server.WithExportType[SeasonInfo](),
		server.WithExportType[SeasonListResponse](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.GET("/api/v1/seasons/{seasonID}",
		h.seasonDetails,
		server.WithExportType[SeasonInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

func (h *SeasonHandler) seasonDetails(w http.ResponseWriter, r *http.Request) {
	seasonID := r.PathValue("seasonID")
	season, err := h.seasonService.GetByID(r.Context(), seasonID)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read season"})
		return
	}

	result := toSeasonInfo(season)
	respond.JSON(w, http.StatusOK, result)
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
		result = append(result, toSeasonInfo(&season))
	}

	respond.JSON(w, http.StatusOK, SeasonListResponse{Seasons: result})
}

func toSeasonInfo(season *entity.Season) SeasonInfo {
	return SeasonInfo{
		ID:     encoders.EncodeUUID(season.ID),
		Number: season.Number,
	}
}
