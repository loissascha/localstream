package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/go-http-server/server"
	"github.com/loissascha/localstream/internal/middleware"
	"github.com/loissascha/localstream/internal/service"
)

type AuthHandler struct {
	s              *server.Server
	authService    *service.AuthService
	authMiddleware *middleware.AuthMiddleware
}

type AuthRequest struct {
	Username string `json:"username"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type AuthUserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type AuthUserIsAdminResponse struct {
	ID      int64 `json:"id"`
	IsAdmin bool  `json:"is_admin"`
}

func NewAuthHandler(s *server.Server, authService *service.AuthService, authMiddleware *middleware.AuthMiddleware) *AuthHandler {
	return &AuthHandler{
		s:              s,
		authService:    authService,
		authMiddleware: authMiddleware,
	}
}

func (h *AuthHandler) RegisterHandlers() {
	h.s.GETI("/api/auth/users/list", h.listUsers, server.WithExportType[AuthUserResponse]())
	h.s.POSTI("/api/auth/register", h.register, server.WithExportType[AuthResponse]())
	h.s.POSTI("/api/auth/login", h.login, server.WithExportType[AuthResponse]())
	h.s.GETI("/api/auth/user/admin", h.isAdmin, server.WithMiddlewares(h.authMiddleware.RequireAuth), server.WithExportType[AuthUserIsAdminResponse]())
}

func (h *AuthHandler) isAdmin(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticatedUserIDFromContext(r)
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	isAdmin, err := h.authService.IsUserAdmin(r.Context(), userID)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respond.JSON(w, http.StatusOK, AuthUserIsAdminResponse{
		ID:      userID,
		IsAdmin: isAdmin,
	})
}

func (h *AuthHandler) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.authService.List(r.Context())
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"error": "couldn't match data. Server error."})
		return
	}

	result := []AuthUserResponse{}
	for _, u := range users {
		result = append(result, AuthUserResponse{
			ID:       u.ID,
			Username: u.Username,
		})
	}

	respond.JSON(w, http.StatusOK, result)
}

func (h *AuthHandler) register(w http.ResponseWriter, r *http.Request) {
	requestBody, err := decodeAuthRequest(r)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	result, err := h.authService.Register(r.Context(), requestBody.Username)
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

	result, err := h.authService.Login(r.Context(), requestBody.Username)
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
