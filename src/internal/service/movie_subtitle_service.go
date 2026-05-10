package service

import (
	"context"
	"fmt"

	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type MovieSubtitleService struct {
	movieSubtitleRepo repository.MovieSubtitleRepository
}

func NewMovieSubtitleService(movieSubtitleRepo repository.MovieSubtitleRepository) *MovieSubtitleService {
	return &MovieSubtitleService{movieSubtitleRepo: movieSubtitleRepo}
}

func (s *MovieSubtitleService) Create(ctx context.Context, subtitle *entity.MovieSubtitle) error {
	if subtitle == nil {
		return fmt.Errorf("movie subtitle is nil")
	}

	if err := s.movieSubtitleRepo.Create(ctx, subtitle); err != nil {
		return fmt.Errorf("create movie subtitle: %w", err)
	}

	return nil
}

func (s *MovieSubtitleService) GetByID(ctx context.Context, id string) (*entity.MovieSubtitle, error) {
	subtitleID, err := encoders.DecodeUUID(id)
	if err != nil {
		return nil, fmt.Errorf("parse subtitle id: %w", err)
	}

	subtitle, err := s.movieSubtitleRepo.GetByID(ctx, subtitleID)
	if err != nil {
		return nil, fmt.Errorf("get movie subtitle by id: %w", err)
	}

	return subtitle, nil
}

func (s *MovieSubtitleService) GetByPath(ctx context.Context, path string) (*entity.MovieSubtitle, error) {
	subtitle, err := s.movieSubtitleRepo.GetByPath(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("get movie subtitle by path: %w", err)
	}

	return subtitle, nil
}

func (s *MovieSubtitleService) ListByMovieID(ctx context.Context, movieID string) ([]entity.MovieSubtitle, error) {
	movieUUID, err := encoders.DecodeUUID(movieID)
	if err != nil {
		return nil, fmt.Errorf("parse movie id: %w", err)
	}

	subtitles, err := s.movieSubtitleRepo.ListByMovieID(ctx, movieUUID)
	if err != nil {
		return nil, fmt.Errorf("list movie subtitles by movie id: %w", err)
	}

	return subtitles, nil
}

func (s *MovieSubtitleService) DeleteByID(ctx context.Context, id string) error {
	subtitleID, err := encoders.DecodeUUID(id)
	if err != nil {
		return fmt.Errorf("parse subtitle id: %w", err)
	}

	if err := s.movieSubtitleRepo.DeleteByID(ctx, subtitleID); err != nil {
		return fmt.Errorf("delete movie subtitle by id: %w", err)
	}

	return nil
}
