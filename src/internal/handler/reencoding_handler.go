package handler

import (
	"net/http"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/media"
	"github.com/loissascha/localstream/internal/middleware"
)

type ReencodeStatus struct {
	Id       string
	Started  bool
	Finished bool
	Failed   bool
	Error    string
}

type ReencodeResponse struct {
	List []ReencodeStatus
}

type ReencodingHandler struct {
	s                 *server.Server
	authMiddleware    *middleware.AuthMiddleware
	reencodingService *media.ReencodingService
}

func NewReencodingHandler(
	s *server.Server,
	authM *middleware.AuthMiddleware,
	reencodingService *media.ReencodingService,
) *ReencodingHandler {
	return &ReencodingHandler{
		s:                 s,
		authMiddleware:    authM,
		reencodingService: reencodingService,
	}
}

func (h *ReencodingHandler) RegisterRoutes() {
	h.s.GETI("/api/v1/movies/reencodings",
		h.listMovies,
		server.WithExportType[ReencodeResponse](),
		server.WithExportType[ReencodeStatus](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.POSTI("/api/v1/movies/reencoding",
		h.startMovie,
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

func (h *ReencodingHandler) startMovie(w http.ResponseWriter, r *http.Request) {
}

func (h *ReencodingHandler) listMovies(w http.ResponseWriter, r *http.Request) {
	queue := h.reencodingService.MovieQueue()
	result := make([]ReencodeStatus, 0, len(queue))

	for _, q := range queue {
		result = append(result, ReencodeStatus{
			Id:       encoders.EncodeUUID(q.Id),
			Started:  q.Started,
			Finished: q.Finished,
			Failed:   q.Failed,
			Error:    q.Error,
		})
	}

	respond.JSON(w, http.StatusOK, ReencodeResponse{List: result})
}
