package handler

import (
	"net/http"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/service"
)

type SeasonMetadataHandler struct {
	s                 *server.Server
	authMiddleware    *middleware.AuthMiddleware
	seasonMetaService *service.SeasonMetadataService
}

func NewSeasonMetadataHandler(s *server.Server, authM *middleware.AuthMiddleware, seasonMetaService *service.SeasonMetadataService) *SeasonMetadataHandler {
	return &SeasonMetadataHandler{
		s:                 s,
		authMiddleware:    authM,
		seasonMetaService: seasonMetaService,
	}
}

func (h *SeasonMetadataHandler) RegisterRoutes() {
	h.s.GETI("/api/v1/season/metadata/by-season/{seasonID}",
		h.getBySeasonID,
		server.WithExportType[SeasonMetadataInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.GETI("/api/v1/season/metadata/{showID}",
		h.listByShow,
		server.WithExportType[SeasonMetadataInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

func (h *SeasonMetadataHandler) getBySeasonID(w http.ResponseWriter, r *http.Request) {
	seasonID := r.PathValue("seasonID")
	metadata, err := h.seasonMetaService.GetBySeasonID(r.Context(), seasonID)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if metadata == nil {
		respond.JSON(w, http.StatusNotFound, map[string]string{"error": "season metadata not found"})
		return
	}

	respond.JSON(w, http.StatusOK, toSeasonMetadataInfo(metadata))
}

func (h *SeasonMetadataHandler) listByShow(w http.ResponseWriter, r *http.Request) {
	showID := r.PathValue("showID")
	metadata, err := h.seasonMetaService.GetByShowID(r.Context(), showID)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	result := []SeasonMetadataInfo{}
	for _, m := range metadata {
		result = append(result, toSeasonMetadataInfo(&m))
	}

	respond.JSON(w, http.StatusOK, result)
}
