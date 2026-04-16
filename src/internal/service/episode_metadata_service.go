package service

import (
	"context"
	"fmt"

	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type EpisodeMetadataService struct {
	episodeMetadataRepo repository.EpisodeMetadataRepository
}

func NewEpisodeMetadataService(episodeMetadataRepo repository.EpisodeMetadataRepository) *EpisodeMetadataService {
	return &EpisodeMetadataService{episodeMetadataRepo: episodeMetadataRepo}
}

func (s *EpisodeMetadataService) Create(ctx context.Context, episodeID string, metadata *entity.EpisodeMetadata) error {
	episodeUUID, err := encoders.DecodeUUID(episodeID)
	if err != nil {
		return fmt.Errorf("parse episode id: %w", err)
	}

	metadata.EpisodeID = episodeUUID

	if err := s.episodeMetadataRepo.Create(ctx, metadata); err != nil {
		return fmt.Errorf("create episode metadata: %w", err)
	}

	return nil
}

func (s *EpisodeMetadataService) GetByEpisodeID(ctx context.Context, episodeID string) (*entity.EpisodeMetadata, error) {
	episodeUUID, err := encoders.DecodeUUID(episodeID)
	if err != nil {
		return nil, fmt.Errorf("parse episode id: %w", err)
	}

	metadata, err := s.episodeMetadataRepo.GetByEpisodeID(ctx, episodeUUID)
	if err != nil {
		return nil, fmt.Errorf("get episode metadata by episode id: %w", err)
	}

	return metadata, nil
}

func (s *EpisodeMetadataService) GetByShowID(ctx context.Context, showID string) ([]entity.EpisodeMetadata, error) {
	showUUID, err := encoders.DecodeUUID(showID)
	if err != nil {
		return nil, fmt.Errorf("parse show id: %w", err)
	}

	metadata, err := s.episodeMetadataRepo.GetByShowID(ctx, showUUID)
	if err != nil {
		return nil, fmt.Errorf("get episode metadata by show id: %w", err)
	}

	return metadata, nil
}

func (s *EpisodeMetadataService) GetByShowIDAndSeasonNumberAndEpisodeNumber(ctx context.Context, showID string, seasonNumber int, episodeNumber int) (*entity.EpisodeMetadata, error) {
	showUUID, err := encoders.DecodeUUID(showID)
	if err != nil {
		return nil, fmt.Errorf("parse show id: %w", err)
	}

	metadata, err := s.episodeMetadataRepo.GetByShowIDAndSeasonNumberAndEpisodeNumber(ctx, showUUID, seasonNumber, episodeNumber)
	if err != nil {
		return nil, fmt.Errorf("get episode metadata by show id, season number and episode number: %w", err)
	}

	return metadata, nil
}
