package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

var (
	ErrInvalidAuthInput   = errors.New("invalid auth input")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

const defaultJWTSecret = "change-me-in-production"

type AuthResult struct {
	Token string
}

type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret []byte
	tokenTTL  time.Duration
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
	secret := strings.TrimSpace(jwtSecret)
	if secret == "" {
		secret = defaultJWTSecret
	}

	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(secret),
		tokenTTL:  12 * 30 * 24 * time.Hour,
	}
}

func (s *AuthService) IsUserAdmin(ctx context.Context, userID int64) (bool, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}

	return user.IsAdmin, nil
}

func (s *AuthService) Register(ctx context.Context, username string) (*AuthResult, error) {
	username, err := normalizeCredentials(username)
	if err != nil {
		return nil, err
	}

	_, err = s.userRepo.GetByUsername(ctx, username)
	if err == nil {
		return nil, ErrUsernameTaken
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, fmt.Errorf("check existing user: %w", err)
	}

	// hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if err != nil {
	// 	return nil, fmt.Errorf("hash password: %w", err)
	// }

	user := &entity.User{
		Username: username,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		if isUniqueViolation(err) {
			return nil, ErrUsernameTaken
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResult{Token: token}, nil
}

func (s *AuthService) Login(ctx context.Context, username string) (*AuthResult, error) {
	username, err := normalizeCredentials(username)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("get user by username: %w", err)
	}

	// if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
	// 	return nil, ErrInvalidCredentials
	// }

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResult{Token: token}, nil
}

func (s *AuthService) List(ctx context.Context) ([]entity.User, error) {
	return s.userRepo.List(ctx)
}

func (s *AuthService) ValidateToken(tokenString string) (int64, error) {
	trimmedToken := strings.TrimSpace(tokenString)
	if trimmedToken == "" {
		return 0, ErrInvalidToken
	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(trimmedToken, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}

		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, ErrInvalidToken
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil || userID <= 0 {
		return 0, ErrInvalidToken
	}

	return userID, nil
}

func (s *AuthService) generateToken(user *entity.User) (string, error) {
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(user.ID, 10),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenTTL)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("sign jwt token: %w", err)
	}

	return tokenString, nil
}

func normalizeCredentials(username string) (string, error) {
	trimmedUsername := strings.TrimSpace(username)
	if trimmedUsername == "" {
		return "", ErrInvalidAuthInput
	}

	return trimmedUsername, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code == "23505"
}
