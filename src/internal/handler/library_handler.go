package handler

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/entity"
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
	h.s.POSTI("/api/admin/libraries/create",
		h.createLibrary,
		server.WithExportType[CreateLibraryRequest](),
		server.WithExportType[CreateLibraryResponse](),
		server.WithMiddlewares(h.authMiddleware.RequireAuthAdmin),
	)
	h.s.POSTI("/api/admin/libraries/update",
		h.updateLibrary,
		server.WithExportType[UpdateLibraryRequest](),
		server.WithExportType[UpdateLibraryResponse](),
		server.WithMiddlewares(h.authMiddleware.RequireAuthAdmin),
	)
}

type UpdateLibraryRequest struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	LibraryType entity.LibraryType `json:"type"`
	Path        string             `json:"path"`
}

type UpdateLibraryResponse struct {
	Library LibraryListItem `json:"library"`
}

func (h *LibraryHandler) updateLibrary(w http.ResponseWriter, r *http.Request) {
	requestBody, err := decodeUpdateLibraryRequest(r)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	lib, err := h.libraryService.Update(r.Context(), requestBody.ID, requestBody.Name, requestBody.Path, requestBody.LibraryType)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response := UpdateLibraryResponse{
		Library: toLibraryListItem(lib),
	}

	respond.JSON(w, http.StatusOK, response)
}

type CreateLibraryRequest struct {
	Name        string             `json:"name"`
	LibraryType entity.LibraryType `json:"type"`
	Path        string             `json:"path"`
}

type CreateLibraryResponse struct {
	Library LibraryListItem `json:"library"`
}

func (h *LibraryHandler) createLibrary(w http.ResponseWriter, r *http.Request) {
	requestBody, err := decodeCreateLibraryRequest(r)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	lib, err := h.libraryService.Create(r.Context(), requestBody.Name, requestBody.Path, string(requestBody.LibraryType))
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response := CreateLibraryResponse{
		Library: toLibraryListItem(lib),
	}

	respond.JSON(w, http.StatusOK, response)
}

func (h *LibraryHandler) listLibraries(w http.ResponseWriter, r *http.Request) {
	libraries, err := h.libraryService.List(r.Context())
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read libraries"})
		return
	}

	result := []LibraryListItem{}
	for _, l := range libraries {
		result = append(result, toLibraryListItem(&l))
	}

	sort.Slice(result, func(i, j int) bool {
		return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
	})

	respond.JSON(w, http.StatusOK, LibraryListResponse{Libraries: result})
}

func decodeUpdateLibraryRequest(r *http.Request) (*UpdateLibraryRequest, error) {
	defer r.Body.Close()

	var requestBody UpdateLibraryRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, err
	}

	return &requestBody, nil
}

func decodeCreateLibraryRequest(r *http.Request) (*CreateLibraryRequest, error) {
	defer r.Body.Close()

	var requestBody CreateLibraryRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, err
	}

	return &requestBody, nil
}
