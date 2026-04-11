package handler

import (
	"net/http"
	"strconv"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/service"
)

type EpisodeMetadataHandler struct {
	s                  *server.Server
	authMiddleware     *middleware.AuthMiddleware
	episodeMetaService *service.EpisodeMetadataService
}

func NewEpisodeMetadataHandler(s *server.Server, authM *middleware.AuthMiddleware, episodeMetaService *service.EpisodeMetadataService) *EpisodeMetadataHandler {
	return &EpisodeMetadataHandler{
		s:                  s,
		authMiddleware:     authM,
		episodeMetaService: episodeMetaService,
	}
}

func (h *EpisodeMetadataHandler) RegisterRoutes() {
	h.s.GETI("/api/v1/episode/metadata/{showID}",
		h.listByShow,
		server.WithExportType[EpisodeMetadataInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.GETI("/api/v1/episode/metadata/{showID}/{seasonNumber}/{episodeNumber}",
		h.getByShowAndSeasonAndEpisodeNumber,
		server.WithExportType[EpisodeMetadataInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

func (h *EpisodeMetadataHandler) listByShow(w http.ResponseWriter, r *http.Request) {
	showID := r.PathValue("showID")
	metadata, err := h.episodeMetaService.GetByShowID(r.Context(), showID)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	result := []EpisodeMetadataInfo{}
	for _, m := range metadata {
		result = append(result, toEpisodeMetadataInfo(&m))
	}

	respond.JSON(w, http.StatusOK, result)
}

func (h *EpisodeMetadataHandler) getByShowAndSeasonAndEpisodeNumber(w http.ResponseWriter, r *http.Request) {
	showID := r.PathValue("showID")
	seasonNumberRaw := r.PathValue("seasonNumber")
	seasonNumber, err := strconv.Atoi(seasonNumberRaw)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid season number"})
		return
	}

	episodeNumberRaw := r.PathValue("episodeNumber")
	episodeNumber, err := strconv.Atoi(episodeNumberRaw)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid episode number"})
		return
	}

	metadata, err := h.episodeMetaService.GetByShowIDAndSeasonNumberAndEpisodeNumber(r.Context(), showID, seasonNumber, episodeNumber)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if metadata == nil {
		respond.JSON(w, http.StatusNotFound, map[string]string{"error": "episode metadata not found"})
		return
	}

	respond.JSON(w, http.StatusOK, toEpisodeMetadataInfo(metadata))
}
