package handler

import (
	"net/http"
	"sort"
	"strings"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/service"
)

type SearchHandler struct {
	s              *server.Server
	authMiddleware *middleware.AuthMiddleware
	showSerivce    *service.ShowService
	movieService   *service.MovieService
}

func NewSearchHandler(s *server.Server, authMiddleware *middleware.AuthMiddleware, showService *service.ShowService, movieService *service.MovieService) *SearchHandler {
	return &SearchHandler{
		s:              s,
		authMiddleware: authMiddleware,
		showSerivce:    showService,
		movieService:   movieService,
	}
}

func (h *SearchHandler) RegisterRoutes() {
	h.s.GETI("/api/v1/search",
		h.search,
		server.WithExportType[SearchResponse](),
		server.WithExportType[ShowInfo](),
		server.WithExportType[MovieInfo](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

func (h *SearchHandler) search(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if len(query) < 3 {
		respond.JSON(w, http.StatusOK, SearchResponse{
			Shows:  []ShowInfo{},
			Movies: []MovieInfo{},
		})
		return
	}

	shows, err := h.showSerivce.Search(r.Context(), query)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to search shows"})
		return
	}

	movies, err := h.movieService.Search(r.Context(), query)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to search movies"})
		return
	}

	showResult := make([]ShowInfo, 0, len(shows))
	for _, show := range shows {
		showResult = append(showResult, toShowInfo(&show))
	}

	movieResult := make([]MovieInfo, 0, len(movies))
	for _, movie := range movies {
		movieResult = append(movieResult, toMovieInfo(&movie))
	}

	sort.Slice(showResult, func(i, j int) bool {
		leftLower := strings.ToLower(showResult[i].Name)
		rightLower := strings.ToLower(showResult[j].Name)
		if leftLower == rightLower {
			return shows[i].Year < shows[j].Year
		}

		return leftLower < rightLower
	})

	sort.Slice(movieResult, func(i, j int) bool {
		leftLower := strings.ToLower(movieResult[i].Name)
		rightLower := strings.ToLower(movieResult[j].Name)
		if leftLower == rightLower {
			return movies[i].Year < movies[j].Year
		}

		return leftLower < rightLower
	})

	respond.JSON(w, http.StatusOK, SearchResponse{
		Shows:  showResult,
		Movies: movieResult,
	})
}
