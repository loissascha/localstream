package service

import (
	"context"
	"fmt"

	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type MovieMetadataService struct {
	movieRepo         repository.MovieRepository
	movieMetadataRepo repository.MovieMetadataRepository
}

func NewMovieMetadataService(movieMetadataRepo repository.MovieMetadataRepository, movieRepo repository.MovieRepository) *MovieMetadataService {
	return &MovieMetadataService{
		movieMetadataRepo: movieMetadataRepo,
		movieRepo:         movieRepo,
	}
}

func (s *MovieMetadataService) Create(ctx context.Context, movieID string, metadata *entity.MovieMetadata) error {
	movieUUID, err := encoders.DecodeUUID(movieID)
	if err != nil {
		return fmt.Errorf("parse movie id: %w", err)
	}

	metadata.MovieID = movieUUID

	if err := s.movieMetadataRepo.Create(ctx, metadata); err != nil {
		return fmt.Errorf("create movie metadata: %w", err)
	}

	return nil
}

func (s *MovieMetadataService) GetByMovieID(ctx context.Context, movieID string) ([]entity.MovieMetadata, error) {
	movieUUID, err := encoders.DecodeUUID(movieID)
	if err != nil {
		return nil, fmt.Errorf("parse movie id: %w", err)
	}

	metadata, err := s.movieMetadataRepo.GetByMovieID(ctx, movieUUID)
	if err != nil {
		return nil, fmt.Errorf("get movie metadata by movie id: %w", err)
	}

	return metadata, nil
}

func (s *MovieMetadataService) SetPrimaryForMovieID(ctx context.Context, movieID string, id string) error {
	uuid, err := encoders.DecodeUUID(id)
	if err != nil {
		return fmt.Errorf("parse id: %w", err)
	}

	movieUUID, err := encoders.DecodeUUID(movieID)
	if err != nil {
		return fmt.Errorf("parse movie id: %w", err)
	}

	metadata, err := s.movieMetadataRepo.GetByMovieID(ctx, movieUUID)
	if err != nil {
		return fmt.Errorf("get movie metadata by movie id: %w", err)
	}

	targetFetchSource := entity.FetchSourceNone
	for _, m := range metadata {
		if m.ID != uuid {
			err := s.movieMetadataRepo.DeleteOne(ctx, m.ID)
			if err != nil {
				return err
			}
			continue
		}
		targetFetchSource = m.FetchSource
	}

	err = s.movieRepo.UpdateFetchSource(ctx, movieUUID, targetFetchSource)
	if err != nil {
		return err
	}

	return nil
}
