package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/helper"
	"github.com/loissascha/localstream/internal/provider"
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

func (self *SeasonMetadataService) CreateSeasonMetadata(ctx context.Context, season *entity.Season, metadata *provider.SeasonMetadata) error {
	mid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	mediumImage := ""
	originalImage := ""
	if metadata.Image != nil {
		mediumImage, err = helper.DownloadImageAndGetStaticPath(
			metadata.Image.Medium,
			helper.GetShowImagePath(season.ShowID),
			fmt.Sprintf("med_SE_%s", mid.String()),
		)
		if err != nil {
			return err
		}
		originalImage, err = helper.DownloadImageAndGetStaticPath(
			metadata.Image.Original,
			helper.GetShowImagePath(season.ShowID),
			fmt.Sprintf("org_%s", mid.String()),
		)
		if err != nil {
			return err
		}
	}
	m := entity.SeasonMetadata{
		ID:               mid,
		SeasonID:         season.ID,
		Url:              metadata.Url,
		Number:           metadata.Number,
		Summary:          metadata.Summary,
		PremiereDate:     metadata.PremiereDate,
		MediumImageUrl:   mediumImage,
		OriginalImageUrl: originalImage,
		FetchSource:      entity.FetchSourceTVMaze,
		FetchID:          metadata.ID,
	}
	err = self.seasonMetadataRepo.Create(ctx, &m)
	if err != nil {
		return err
	}
	return nil
}
