package service

import (
	"context"
	"fmt"

	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/provider"
	"github.com/loissascha/localstream/internal/repository"
)

type ShowMetadataService struct {
	showRepo             repository.ShowRepository
	showMetadataRepo     repository.ShowMetadataRepository
	showMetadataProvider provider.TVMetadataProvider
}

func NewShowMetadataService(showMetadataRepo repository.ShowMetadataRepository, showRepo repository.ShowRepository, showMetadataProvider provider.TVMetadataProvider) *ShowMetadataService {
	return &ShowMetadataService{
		showMetadataRepo:     showMetadataRepo,
		showRepo:             showRepo,
		showMetadataProvider: showMetadataProvider,
	}
}

func (s *ShowMetadataService) Search(ctx context.Context, searchQuery string) ([]provider.ShowSearchResult, error) {
	return s.showMetadataProvider.SearchShow(searchQuery, 0)
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

func (s *ShowMetadataService) SetPrimaryForShowID(ctx context.Context, showID string, id string) error {
	uuid, err := encoders.DecodeUUID(id)
	if err != nil {
		return fmt.Errorf("parse id: %w", err)
	}

	showUUID, err := encoders.DecodeUUID(showID)
	if err != nil {
		return fmt.Errorf("parse show id: %w", err)
	}

	metadata, err := s.showMetadataRepo.GetByShowID(ctx, showUUID)
	if err != nil {
		return fmt.Errorf("get show metadata by show id: %w", err)
	}

	targetFetchSource := entity.FetchSourceNone
	for _, m := range metadata {
		if m.ID != uuid {
			err := s.showMetadataRepo.DeleteOne(ctx, m.ID)
			if err != nil {
				return err
			}
			continue
		}
		targetFetchSource = m.FetchSource
	}

	err = s.showRepo.UpdateFetchSource(ctx, showUUID, targetFetchSource)
	if err != nil {
		return err
	}

	return nil
}
