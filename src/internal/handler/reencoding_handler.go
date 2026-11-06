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
	Id       string `json:"id"`
	Started  bool   `json:"started"`
	Finished bool   `json:"finished"`
	Failed   bool   `json:"failed"`
	Error    string `json:"error"`
}

type ReencodeResponse struct {
	List []ReencodeStatus `json:"list"`
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
	h.s.GETI("/api/v1/reencoding/movies",
		h.listMovies,
		server.WithExportType[ReencodeResponse](),
		server.WithExportType[ReencodeStatus](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)

	h.s.POSTI("/api/v1/reencoding/movie/{movieID}",
		h.startMovie,
		server.WithExportType[ReencodeStatus](),
		server.WithMiddlewares(h.authMiddleware.RequireAuth),
	)
}

func (h *ReencodingHandler) startMovie(w http.ResponseWriter, r *http.Request) {
	movieId := r.PathValue("movieID")

	mid, err := encoders.DecodeUUID(movieId)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to convert id: " + err.Error()})
		return
	}

	q := h.reencodingService.AddMovieToReencodeQueue(mid)
	result := ReencodeStatus{
		Id:       encoders.EncodeUUID(q.Id),
		Started:  q.Started,
		Finished: q.Finished,
		Failed:   q.Failed,
		Error:    q.Error,
	}
	respond.JSON(w, http.StatusOK, result)
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
