package service

import (
	"context"
	"fmt"

	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type SeasonMetadataService struct {
	seasonMetadataRepo repository.SeasonMetadataRepository
}

func NewSeasonMetadataService(seasonMetadataRepo repository.SeasonMetadataRepository) *SeasonMetadataService {
	return &SeasonMetadataService{seasonMetadataRepo: seasonMetadataRepo}
}

func (s *SeasonMetadataService) Create(ctx context.Context, seasonID string, metadata *entity.SeasonMetadata) error {
	seasonUUID, err := encoders.DecodeUUID(seasonID)
	if err != nil {
		return fmt.Errorf("parse season id: %w", err)
	}

	metadata.SeasonID = seasonUUID

	if err := s.seasonMetadataRepo.Create(ctx, metadata); err != nil {
		return fmt.Errorf("create season metadata: %w", err)
	}

	return nil
}

func (s *SeasonMetadataService) GetBySeasonID(ctx context.Context, seasonID string) (*entity.SeasonMetadata, error) {
	seasonUUID, err := encoders.DecodeUUID(seasonID)
	if err != nil {
		return nil, fmt.Errorf("parse season id: %w", err)
	}

	metadata, err := s.seasonMetadataRepo.GetBySeasonID(ctx, seasonUUID)
	if err != nil {
		return nil, fmt.Errorf("get season metadata by season id: %w", err)
	}

	return metadata, nil
}

func (s *SeasonMetadataService) GetByShowID(ctx context.Context, showID string) ([]entity.SeasonMetadata, error) {
	showUUID, err := encoders.DecodeUUID(showID)
	if err != nil {
		return nil, fmt.Errorf("parse show id: %w", err)
	}

	metadata, err := s.seasonMetadataRepo.GetByShowID(ctx, showUUID)
	if err != nil {
		return nil, fmt.Errorf("get season metadata by show id: %w", err)
	}

	return metadata, nil
}
