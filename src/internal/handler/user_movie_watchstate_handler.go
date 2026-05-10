package handler

import (
	"errors"
	"net/http"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/service"
)

type UserMovieWatchstateHandler struct {
	s                          *server.Server
	authMiddleware             *middleware.AuthMiddleware
	userMovieWatchstateService *service.UserMovieWatchstateService
	movieService               *service.MovieService
}

func NewUserMovieWatchstateHandler(
	s *server.Server,
	authM *middleware.AuthMiddleware,
	userWatchstateS *service.UserMovieWatchstateService,
	movieS *service.MovieService,
) *UserMovieWatchstateHandler {
	return &UserMovieWatchstateHandler{
		s:                          s,
		authMiddleware:             authM,
		userMovieWatchstateService: userWatchstateS,
		movieService:               movieS,
	}
}

func (h *UserMovieWatchstateHandler) RegisterHandlers() {
	h.s.POST("/api/v1/movie/watchstate",
		h.saveWatchstate,
		server.WithExportType[SaveMovieWatchstateRequest](),
		server.WithExportType[WatchstateInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
		server.WithDescription("Save the current watchstate for movie"),
	)

	h.s.GET("/api/v1/watchstate/movie/{movieID}",
		h.getWatchstateByMovieID,
		server.WithExportType[WatchstateInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
		server.WithDescription("Get the last watchstate of the movie"),
	)

	h.s.GET("/api/v1/watchstate/movie/latest",
		h.listLatestWatchstates,
		server.WithExportType[WatchstateMoviesListResponse](),
		server.WithExportType[WatchstateMovieResponse](),
		server.WithExportType[MovieInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.POST("/api/v1/watchstate/movie/{movieID}/finished",
		h.setMovieWatchstateFinished,
		server.WithExportType[WatchstateMovieResponse](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
		server.WithDescription("Sets the movie to watched"),
	)

	h.s.DELETEI("/api/v1/watchstate/movie/{movieID}/delete",
		h.deleteMovieWatchstate,
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
		server.WithDescription("deletes the watchstate of a movie"),
	)
}

func (h *UserMovieWatchstateHandler) setMovieWatchstateFinished(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	userID, ok := authenticatedUserIDFromContext(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	watchstate, err := h.userMovieWatchstateService.Save(r.Context(), userID, service.SaveMovieWatchstateInput{
		MovieID:  movieID,
		Position: 0,
		Duration: 0,
		Finished: true,
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

func (h *UserMovieWatchstateHandler) deleteMovieWatchstate(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	userID, ok := authenticatedUserIDFromContext(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	err := h.userMovieWatchstateService.DeleteByMovieID(r.Context(), userID, movieID)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	respond.JSON(w, http.StatusOK, true)
}

func (h *UserMovieWatchstateHandler) listLatestWatchstates(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticatedUserIDFromContext(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	watchstates, err := h.userMovieWatchstateService.ListByUserID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidWatchstateInput) {
			respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
			return
		}

		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read watchstates"})
		return
	}

	response := make([]WatchstateMovieResponse, 0, len(watchstates))
	for _, watchstate := range watchstates {

		if watchstate.Finished {
			continue
		}

		movieId := encoders.EncodeUUID(watchstate.MovieID)
		movie, err := h.movieService.GetByIDWithMetadata(r.Context(), movieId, userID)
		if err != nil {
			respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read watchstates: " + err.Error()})
			return
		}

		movieInfo := toMovieInfo(movie)
		response = append(response, toWatchstateMovieResponse(watchstate, movieInfo))
	}

	respond.JSON(w, http.StatusOK, WatchstateMoviesListResponse{Watchstates: response})
}

func (h *UserMovieWatchstateHandler) getWatchstateByMovieID(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticatedUserIDFromContext(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	movieID := r.PathValue("movieID")
	watchstate, err := h.userMovieWatchstateService.GetByMovieID(r.Context(), userID, movieID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidWatchstateInput) {
			respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid movie id"})
			return
		}

		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read watchstate"})
		return
	}

	if watchstate == nil {
		respond.JSON(w, http.StatusNotFound, map[string]string{"error": "watchstate not found"})
		return
	}

	respond.JSON(w, http.StatusOK, toWatchstateInfoMovie(watchstate))
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
		if percent >= 90 {
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
