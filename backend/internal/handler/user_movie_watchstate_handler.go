package handler

import (
	"errors"
	"net/http"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/service"
)

type UserMovieWatchstateHandler struct {
	s                          *server.Server
	authMiddleware             *middleware.AuthMiddleware
	userMovieWatchstateService *service.UserMovieWatchstateService
	movieService               *service.MovieService
}

func NewUserMovieWatchstateHandler(s *server.Server, authM *middleware.AuthMiddleware, userWatchstateS *service.UserMovieWatchstateService, movieS *service.MovieService) *UserMovieWatchstateHandler {
	return &UserMovieWatchstateHandler{
		s:                          s,
		authMiddleware:             authM,
		userMovieWatchstateService: userWatchstateS,
		movieService:               movieS,
	}
}

func (h *UserMovieWatchstateHandler) RegisterHandlers() {
	h.s.POST("/api/watchstate",
		h.saveWatchstate,
		server.WithExportType[SaveMovieWatchstateRequest](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
		server.WithDescription("Save the current watchstate for movie"),
	)
}

type SaveMovieWatchstateRequest struct {
	MovieID  string  `json:"movie_id"`
	Position float64 `json:"position"`
	Duration float64 `json:"duration"`
	Finished bool    `json:"finished"`
}

func (h *UserMovieWatchstateHandler) saveWatchstate(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticatedUserIDFromContext(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	requestBody, err := decodeSaveMovieWatchstateRequest(r)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	// set finish when almost done
	finished := requestBody.Finished
	if !finished {
		percent := (100 / requestBody.Duration) * requestBody.Position
		if percent >= 95 {
			finished = true
		}
	}

	watchstate, err := h.userMovieWatchstateService.Save(r.Context(), userID, service.SaveMovieWatchstateInput{
		MovieID:  requestBody.MovieID,
		Position: requestBody.Position,
		Duration: requestBody.Duration,
		Finished: finished,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidWatchstateInput) {
			respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid watchstate input"})
			return
		}

		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to save watchstate"})
		return
	}

	respond.JSON(w, http.StatusOK, toWatchstateInfoMovie(watchstate))
}
