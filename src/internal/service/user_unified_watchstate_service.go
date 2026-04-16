package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

var ErrInvalidUnifiedWatchstateInput = errors.New("invalid unified watchstate input")

type UserUnifiedWatchstateService struct {
	watchstateRepo repository.UserUnifiedWatchstateRepository
}

func NewUserUnifiedWatchstateService(watchstateRepo repository.UserUnifiedWatchstateRepository) *UserUnifiedWatchstateService {
	return &UserUnifiedWatchstateService{watchstateRepo: watchstateRepo}
}

func (s *UserUnifiedWatchstateService) ListByUserID(ctx context.Context, userID int64) ([]entity.UserUnifiedWatchstate, error) {
	if userID <= 0 {
		return nil, ErrInvalidUnifiedWatchstateInput
	}

	watchstates, err := s.watchstateRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list unified watchstates by user id: %w", err)
	}

	return watchstates, nil
}
