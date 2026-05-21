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

type EpisodeMetadataService struct {
	episodeMetadataRepo repository.EpisodeMetadataRepository
	seasonRepo          repository.SeasonRepository
}

func NewEpisodeMetadataService(
	episodeMetadataRepo repository.EpisodeMetadataRepository,
	seasonRepo repository.SeasonRepository,
) *EpisodeMetadataService {
	return &EpisodeMetadataService{
		episodeMetadataRepo: episodeMetadataRepo,
		seasonRepo:          seasonRepo,
	}
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

func (self *EpisodeMetadataService) CreateEpisodeMetadata(ctx context.Context, episode *entity.Episode, metadata *provider.EpisodeMetadata) error {
	mid, err := uuid.NewV7()
	if err != nil {
		return err
	}

	season, err := self.seasonRepo.GetByID(ctx, episode.SeasonID)
	if err != nil {
		return err
	}

	mediumImage := ""
	originalImage := ""
	if metadata.Image != nil {
		mediumImage, err = helper.DownloadImageAndGetStaticPath(
			metadata.Image.Medium,
			helper.GetShowImagePath(season.ShowID),
			fmt.Sprintf("med_E_%s", mid.String()),
		)
		if err != nil {
			return err
		}
		originalImage, err = helper.DownloadImageAndGetStaticPath(
			metadata.Image.Original,
			helper.GetShowImagePath(season.ShowID),
			fmt.Sprintf("org_E_%s", mid.String()),
		)
		if err != nil {
			return err
		}
	}
	m := entity.EpisodeMetadata{
		ID:               mid,
		EpisodeID:        episode.ID,
		Url:              metadata.Url,
		Name:             metadata.Name,
		Number:           metadata.Number,
		Summary:          metadata.Summary,
		MediumImageUrl:   mediumImage,
		OriginalImageUrl: originalImage,
		FetchSource:      entity.FetchSourceTVMaze,
		FetchID:          metadata.ID,
	}
	err = self.episodeMetadataRepo.Create(ctx, &m)
	if err != nil {
		return err
	}
	return nil
}
