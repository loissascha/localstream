package homehandler

import (
	"net/http"

	"github.com/loissascha/go-http-server/server"
)

type HomeHandler struct {
	s *server.Server
}

func New(s *server.Server) *HomeHandler {
	return &HomeHandler{
		s: s,
	}
}

func (h *HomeHandler) RegisterHandlers() {
	h.s.GET("/", h.homeRoute)
}

func (h *HomeHandler) homeRoute(w http.ResponseWriter, r *http.Request) {
}
