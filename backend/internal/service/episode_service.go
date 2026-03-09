package service

import (
	"context"
	"fmt"

	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type EpisodeService struct {
	episodeRepo repository.EpisodeRepository
}

func NewEpisodeService(episodeRepo repository.EpisodeRepository) *EpisodeService {
	return &EpisodeService{episodeRepo: episodeRepo}
}

func (s *EpisodeService) ListBySeasonID(ctx context.Context, seasonId string) ([]entity.Episode, error) {
	seasonUUID, err := encoders.DecodeUUID(seasonId)
	if err != nil {
		return nil, fmt.Errorf("parse season id: %w", err)
	}

	episodes, err := s.episodeRepo.ListBySeasonID(ctx, seasonUUID)
	if err != nil {
		return nil, fmt.Errorf("list episodes by season id: %w", err)
	}

	return episodes, nil
}

func (s *EpisodeService) GetByID(ctx context.Context, episodeId string) (*entity.Episode, error) {
	id, err := encoders.DecodeUUID(episodeId)
	if err != nil {
		return nil, err
	}

	return s.episodeRepo.GetByID(ctx, id)
}
