package service

import (
	"context"
	"fmt"

	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type ShowMetadataService struct {
	showMetadataRepo repository.ShowMetadataRepository
}

func NewShowMetadataService(showMetadataRepo repository.ShowMetadataRepository) *ShowMetadataService {
	return &ShowMetadataService{showMetadataRepo: showMetadataRepo}
}

func (s *ShowMetadataService) Create(ctx context.Context, showID string, metadata *entity.ShowMetadata) error {
	showUUID, err := encoders.DecodeUUID(showID)
	if err != nil {
		return fmt.Errorf("parse show id: %w", err)
	}

	metadata.ShowID = showUUID

	if err := s.showMetadataRepo.Create(ctx, metadata); err != nil {
		return fmt.Errorf("create show metadata: %w", err)
	}

	return nil
}

func (s *ShowMetadataService) GetByShowID(ctx context.Context, showID string) ([]entity.ShowMetadata, error) {
	showUUID, err := encoders.DecodeUUID(showID)
	if err != nil {
		return nil, fmt.Errorf("parse show id: %w", err)
	}

	metadata, err := s.showMetadataRepo.GetByShowID(ctx, showUUID)
	if err != nil {
		return nil, fmt.Errorf("get show metadata by show id: %w", err)
	}

	return metadata, nil
}
