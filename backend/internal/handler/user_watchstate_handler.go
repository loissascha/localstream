package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/service"
)

type UserWatchstateHandler struct {
	s                     *server.Server
	authMiddleware        *middleware.AuthMiddleware
	userWatchstateService *service.UserWatchstateService
	showService           *service.ShowService
	seasonService         *service.SeasonService
	episodeService        *service.EpisodeService
}

func NewUserWatchstateHandler(
	s *server.Server,
	authMiddleware *middleware.AuthMiddleware,
	userWatchstateService *service.UserWatchstateService,
	showService *service.ShowService,
	seasonService *service.SeasonService,
	episodeService *service.EpisodeService,
) *UserWatchstateHandler {
	return &UserWatchstateHandler{
		s:                     s,
		authMiddleware:        authMiddleware,
		userWatchstateService: userWatchstateService,
		showService:           showService,
		seasonService:         seasonService,
		episodeService:        episodeService,
	}
}

func (h *UserWatchstateHandler) RegisterRoutes() {
	h.s.POST("/api/watchstate",
		h.saveWatchstate,
		server.WithExportType[SaveWatchstateRequest](),
		server.WithExportType[WatchstateResponse](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.GET("/api/watchstate/episode/{episodeID}",
		h.getWatchstateByEpisodeID,
		server.WithExportType[WatchstateResponse](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.GET("/api/watchstate/latest/shows",
		h.listLatestWatchstatesByShow,
		server.WithExportType[WatchstateResponse](),
		server.WithExportType[WatchstateListResponse](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

type SaveWatchstateRequest struct {
	ShowID    string  `json:"show_id"`
	SeasonID  string  `json:"season_id"`
	EpisodeID string  `json:"episode_id"`
	Position  float64 `json:"position"`
	Duration  float64 `json:"duration"`
	Finished  bool    `json:"finished"`
}

type WatchstateResponse struct {
	ID          string      `json:"id"`
	ShowID      string      `json:"show_id"`
	ShowInfo    ShowInfo    `json:"show_info"`
	SeasonID    string      `json:"season_id"`
	SeasonInfo  SeasonInfo  `json:"season_info"`
	EpisodeID   string      `json:"episode_id"`
	EpisodeInfo EpisodeInfo `json:"episode_info"`
	Position    float64     `json:"position"`
	Duration    float64     `json:"duration"`
	Finished    bool        `json:"finished"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type WatchstateListResponse struct {
	Watchstates []WatchstateResponse `json:"watchstates"`
}

func (h *UserWatchstateHandler) saveWatchstate(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticatedUserIDFromContext(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	requestBody, err := decodeSaveWatchstateRequest(r)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	watchstate, err := h.userWatchstateService.Save(r.Context(), userID, service.SaveWatchstateInput{
		ShowID:    requestBody.ShowID,
		SeasonID:  requestBody.SeasonID,
		EpisodeID: requestBody.EpisodeID,
		Position:  requestBody.Position,
		Duration:  requestBody.Duration,
		Finished:  requestBody.Finished,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidWatchstateInput) {
			respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid watchstate input"})
			return
		}

		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to save watchstate"})
		return
	}

	respond.JSON(w, http.StatusOK, toWatchstateResponse(*watchstate))
}

func (h *UserWatchstateHandler) getWatchstateByEpisodeID(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticatedUserIDFromContext(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	episodeID := r.PathValue("episodeID")
	watchstate, err := h.userWatchstateService.GetByEpisodeID(r.Context(), userID, episodeID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidWatchstateInput) {
			respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid episode id"})
			return
		}

		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read watchstate"})
		return
	}

	if watchstate == nil {
		respond.JSON(w, http.StatusNotFound, map[string]string{"error": "watchstate not found"})
		return
	}

	respond.JSON(w, http.StatusOK, toWatchstateResponse(*watchstate))
}

func (h *UserWatchstateHandler) listLatestWatchstatesByShow(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticatedUserIDFromContext(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	watchstates, err := h.userWatchstateService.ListLatestByShow(r.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidWatchstateInput) {
			respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
			return
		}

		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read watchstates"})
		return
	}

	response := make([]WatchstateResponse, 0, len(watchstates))
	for _, watchstate := range watchstates {
		showId := encoders.EncodeUUID(watchstate.ShowID)
		seasonId := encoders.EncodeUUID(watchstate.SeasonID)
		episodeId := encoders.EncodeUUID(watchstate.EpisodeID)

		show, err := h.showService.GetByID(r.Context(), showId)
		if err != nil {
			respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read watchstates: " + err.Error()})
			return
		}

		season, err := h.seasonService.GetByID(r.Context(), seasonId)
		if err != nil {
			respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read watchstates: " + err.Error()})
			return
		}

		episode, err := h.episodeService.GetByID(r.Context(), episodeId)
		if err != nil {
			respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read watchstates: " + err.Error()})
			return
		}

		response = append(response, toWatchstateResponse(watchstate, toShowInfo(show), toSeasonInfo(season), toEpisodeInfo(episode)))
	}

	respond.JSON(w, http.StatusOK, WatchstateListResponse{Watchstates: response})
}

func decodeSaveWatchstateRequest(r *http.Request) (*SaveWatchstateRequest, error) {
	defer r.Body.Close()

	var requestBody SaveWatchstateRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, err
	}

	return &requestBody, nil
}

func authenticatedUserIDFromContext(r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value(middleware.AuthenticatedUserIDKey).(int64)
	if !ok || userID <= 0 {
		return 0, false
	}

	return userID, true
}
