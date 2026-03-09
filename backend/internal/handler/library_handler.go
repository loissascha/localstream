package handler

import (
	"net/http"
	"sort"
	"strings"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/service"
)

type LibraryHandler struct {
	s              *server.Server
	authMiddleware *middleware.AuthMiddleware
	libraryService *service.LibraryService
}

func NewLibraryHandler(s *server.Server, authMiddleware *middleware.AuthMiddleware, libService *service.LibraryService) *LibraryHandler {
	return &LibraryHandler{
		s:              s,
		authMiddleware: authMiddleware,
		libraryService: libService,
	}
}

func (h *LibraryHandler) RegisterHandlers() {
	h.s.GET("/api/libraries",
		h.listLibraries,
		server.WithExportType[LibraryListItem](),
		server.WithExportType[LibraryListResponse](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

type LibraryListItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	LibraryType string `json:"library_type"`
}

type LibraryListResponse struct {
	Libraries []LibraryListItem `json:"libraries"`
}

func (h *LibraryHandler) listLibraries(w http.ResponseWriter, r *http.Request) {
	libraries, err := h.libraryService.List(r.Context())
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read libraries"})
		return
	}

	result := []LibraryListItem{}
	for _, l := range libraries {
		result = append(result, LibraryListItem{
			ID:          encoders.EncodeUUID(l.ID),
			Name:        l.Name,
			Path:        l.Path,
			LibraryType: string(l.LibraryType),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
	})

	respond.JSON(w, http.StatusOK, LibraryListResponse{Libraries: result})
}
