package handler

import (
	"net/http"

	"github.com/loissascha/go-http-server/respond"
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
	h.s.GETI("/api/v1/show/metadata/{showID}",
		h.listByShow,
		server.WithExportType[ShowMetadataInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.POSTI("/api/v1/show/metadata/{showID}/set/primary/{id}",
		h.setPrimary,
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

func (h *ShowMetadataHandler) setPrimary(w http.ResponseWriter, r *http.Request) {
	showId := r.PathValue("showID")
	id := r.PathValue("id")

	err := h.showMetaService.SetPrimaryForShowID(r.Context(), showId, id)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respond.JSON(w, http.StatusOK, id)
}

func (h *ShowMetadataHandler) listByShow(w http.ResponseWriter, r *http.Request) {
	showId := r.PathValue("showID")
	metadata, err := h.showMetaService.GetByShowID(r.Context(), showId)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	result := []ShowMetadataInfo{}
	for _, m := range metadata {
		result = append(result, toMetadataInfo(&m))
	}

	respond.JSON(w, http.StatusOK, result)
}
