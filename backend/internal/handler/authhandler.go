package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/service"
)

type AuthHandler struct {
	s           *server.Server
	authService *service.AuthService
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func NewAuthHandler(s *server.Server, authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		s:           s,
		authService: authService,
	}
}

func (h *AuthHandler) RegisterHandlers() {
	h.s.POST("/auth/register", h.register, server.WithExportType[AuthResponse]())
	h.s.POST("/auth/reigster", h.register, server.WithExportType[AuthResponse]())
	h.s.POST("/auth/login", h.login, server.WithExportType[AuthResponse]())
}

func (h *AuthHandler) register(w http.ResponseWriter, r *http.Request) {
	requestBody, err := decodeAuthRequest(r)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	result, err := h.authService.Register(r.Context(), requestBody.Username, requestBody.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidAuthInput):
			respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "username and password are required"})
		case errors.Is(err, service.ErrUsernameTaken):
			respond.JSON(w, http.StatusConflict, map[string]string{"error": "username already taken"})
		default:
			respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to register user"})
		}
		return
	}

	respond.JSON(w, http.StatusCreated, AuthResponse{Token: result.Token})
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	requestBody, err := decodeAuthRequest(r)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	result, err := h.authService.Login(r.Context(), requestBody.Username, requestBody.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidAuthInput):
			respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "username and password are required"})
		case errors.Is(err, service.ErrInvalidCredentials):
			respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid username or password"})
		default:
			respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to authenticate user"})
		}
		return
	}

	respond.JSON(w, http.StatusOK, AuthResponse{Token: result.Token})
}

func decodeAuthRequest(r *http.Request) (*AuthRequest, error) {
	defer r.Body.Close()

	var requestBody AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, err
	}

	return &requestBody, nil
}
