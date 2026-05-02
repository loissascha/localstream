package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/provider"
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
		server.WithMiddlewares(h.authMiddleware.RequireAuthAdmin),
	)

	h.s.POSTI("/api/v1/show/metadata/{showID}/set/primary/by-fetchid/{id}",
		h.setPrimaryByFetchID,
		server.WithMiddlewares(h.authMiddleware.RequireAuthAdmin),
	)

	h.s.POSTI("/api/v1/show/metadata/search",
		h.search,
		server.WithExportType[provider.ShowSearchResult](),
		server.WithExportType[provider.ShowMetadata](),
		server.WithExportType[provider.ShowImage](),
		server.WithMiddlewares(h.authMiddleware.RequireAuthAdmin),
	)
}

func (h *ShowMetadataHandler) setPrimaryByFetchID(w http.ResponseWriter, r *http.Request) {
	showId := r.PathValue("showID")
	id := r.PathValue("id")
	idint, err := strconv.Atoi(id)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	err = h.showMetaService.SetPrimaryForShowIDByFetchID(r.Context(), showId, idint)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respond.JSON(w, http.StatusOK, showId)
}

func (h *ShowMetadataHandler) search(w http.ResponseWriter, r *http.Request) {
	searchQuery := strings.TrimSpace(r.URL.Query().Get("q"))
	if searchQuery == "" {
		http.Error(w, "missing search parameter", http.StatusBadRequest)
		return
	}
	result, err := h.showMetaService.Search(r.Context(), searchQuery)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	respond.JSON(w, http.StatusOK, result)
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
		result = append(result, toShowMetadataInfo(&m))
	}

	respond.JSON(w, http.StatusOK, result)
}
