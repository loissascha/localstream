package service

import (
	"context"
	"fmt"

	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type SeasonService struct {
	seasonRepo repository.SeasonRepository
}

func NewSeasonService(seasonRepo repository.SeasonRepository) *SeasonService {
	return &SeasonService{seasonRepo: seasonRepo}
}

func (s *SeasonService) ListByShowID(ctx context.Context, showId string) ([]entity.Season, error) {
	showUUID, err := encoders.DecodeUUID(showId)
	if err != nil {
		return nil, fmt.Errorf("parse show id: %w", err)
	}

	seasons, err := s.seasonRepo.ListByShowID(ctx, showUUID)
	if err != nil {
		return nil, fmt.Errorf("list seasons by show id: %w", err)
	}

	return seasons, nil
}
