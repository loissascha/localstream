package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/loissascha/go-http-server/respond"
	"github.com/loissascha/localstream/internal/service"
)

type contextKey string

const AuthenticatedUserIDKey contextKey = "authenticatedUserID"

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := parseAuthorizationHeader(r.Header.Get("Authorization"))
		if token == "" {
			respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
			return
		}

		userID, err := m.authService.ValidateToken(token)
		if err != nil {
			if errors.Is(err, service.ErrInvalidToken) {
				respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
				return
			}

			respond.JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}

		ctx := context.WithValue(r.Context(), AuthenticatedUserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func parseAuthorizationHeader(headerValue string) string {
	trimmedHeader := strings.TrimSpace(headerValue)
	if trimmedHeader == "" {
		return ""
	}

	parts := strings.Fields(trimmedHeader)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}

	return trimmedHeader
}
