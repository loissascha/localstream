package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/provider"
	"github.com/loissascha/localstream/internal/repository"
)

type ShowMetadataService struct {
	showService          *ShowService
	showRepo             repository.ShowRepository
	showMetadataRepo     repository.ShowMetadataRepository
	showMetadataProvider provider.TVMetadataProvider
}

func NewShowMetadataService(showMetadataRepo repository.ShowMetadataRepository, showRepo repository.ShowRepository, showMetadataProvider provider.TVMetadataProvider, showService *ShowService) *ShowMetadataService {
	return &ShowMetadataService{
		showService:          showService,
		showMetadataRepo:     showMetadataRepo,
		showRepo:             showRepo,
		showMetadataProvider: showMetadataProvider,
	}
}

func (s *ShowMetadataService) SetPrimaryForShowIDByFetchID(ctx context.Context, showID string, id int) error {
	show, err := s.showService.GetByID(ctx, showID)
	if err != nil {
		return err
	}
	if show == nil {
		return fmt.Errorf("show not found")
	}

	// get the data from the provider
	showResult, err := s.showMetadataProvider.GetShowByID(id)
	// if none -> error
	if err != nil {
		return err
	}
	if showResult == nil {
		return fmt.Errorf("show not found (metadata provider)")
	}

	// delete all the current existing metadata for the movie (and reset the fetch state)
	// TODO: make this more performant by adding a bulk delete by movie ID
	metadatas, err := s.showMetadataRepo.GetByShowID(ctx, show.ID)
	if err != nil {
		return err
	}
	for _, m := range metadatas {
		err := s.showMetadataRepo.DeleteOne(ctx, m.ID)
		if err != nil {
			return err
		}
	}

	// create new metadata from the provider result
	err = s.CreateShowMetadata(ctx, showID, showResult)
	if err != nil {
		return err
	}

	// set the fetch result
	err = s.showRepo.UpdateFetchSource(ctx, show.ID, entity.FetchSourceTVMaze)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShowMetadataService) Search(ctx context.Context, searchQuery string) ([]provider.ShowSearchResult, error) {
	return s.showMetadataProvider.SearchShow(searchQuery, 0)
}

func (s *ShowMetadataService) CreateShowMetadata(ctx context.Context, showID string, metadata *provider.ShowMetadata) error {
	uid, err := uuid.NewV7()
	if err != nil {
		return err
	}

	m := entity.ShowMetadata{
		ID:               uid,
		Name:             metadata.Name,
		Url:              metadata.URL,
		Description:      *metadata.Summary,
		MediumImageUrl:   metadata.Image.Medium,
		OriginalImageUrl: metadata.Image.Original,
		FetchID:          metadata.ID,
		FetchSource:      entity.FetchSourceTVMaze,
	}
	return s.Create(ctx, showID, &m)
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
