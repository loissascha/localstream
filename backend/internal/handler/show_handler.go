package handler

import (
	"net/http"
	"sort"
	"strings"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/service"
)

type ShowHandler struct {
	s              *server.Server
	authMiddleware *middleware.AuthMiddleware
	showSerivce    *service.ShowService
}

func NewShowHandler(s *server.Server, authMiddleware *middleware.AuthMiddleware, showService *service.ShowService) *ShowHandler {
	return &ShowHandler{
		s:              s,
		authMiddleware: authMiddleware,
		showSerivce:    showService,
	}
}

func (h *ShowHandler) RegisterRoutes() {
	h.s.GET("/api/shows",
		h.listShows,
		server.WithExportType[ShowInfo](),
		server.WithExportType[ShowListResponse](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
	h.s.GET("/api/show/{id}",
		h.showData,
		server.WithExportType[ShowInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

type ShowInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Year        int    `json:"year"`
	Description string `json:"description"`
}

type ShowListResponse struct {
	Shows []ShowInfo `json:"shows"`
}

func (h *ShowHandler) showData(w http.ResponseWriter, r *http.Request) {
	showId := r.PathValue("id")
	show, err := h.showSerivce.GetByID(r.Context(), showId)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read show"})
		return
	}

	result := toShowInfo(show)
	respond.JSON(w, http.StatusOK, result)
}

func (h *ShowHandler) listShows(w http.ResponseWriter, r *http.Request) {
	shows, err := h.showSerivce.List(r.Context())
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read shows"})
		return
	}

	result := []ShowInfo{}
	for _, l := range shows {
		result = append(result, toShowInfo(&l))
	}

	sort.Slice(result, func(i, j int) bool {
		return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
	})

	respond.JSON(w, http.StatusOK, ShowListResponse{Shows: result})
}

func toShowInfo(show *entity.Show) ShowInfo {
	return ShowInfo{
		ID:          encoders.EncodeUUID(show.ID),
		Name:        show.Name,
		Year:        show.Year,
		Description: show.Description,
	}
}
