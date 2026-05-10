package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

var ErrInvalidMovieWatchstateInput = errors.New("invalid movie watchstate input")

type SaveMovieWatchstateInput struct {
	MovieID  string
	Position float64
	Duration float64
	Finished bool
}

type UserMovieWatchstateService struct {
	watchstateRepo repository.UserMovieWatchstateRepository
}

func NewUserMovieWatchstateService(watchstateRepo repository.UserMovieWatchstateRepository) *UserMovieWatchstateService {
	return &UserMovieWatchstateService{watchstateRepo: watchstateRepo}
}

func (s *UserMovieWatchstateService) Save(ctx context.Context, userID int64, input SaveMovieWatchstateInput) (*entity.UserMovieWatchstate, error) {
	if userID <= 0 {
		return nil, ErrInvalidMovieWatchstateInput
	}
	if input.Position < 0 || input.Duration < 0 {
		return nil, ErrInvalidMovieWatchstateInput
	}

	movieID, err := encoders.DecodeUUID(input.MovieID)
	if err != nil {
		return nil, fmt.Errorf("decode movie id: %w", ErrInvalidMovieWatchstateInput)
	}

	watchstate := &entity.UserMovieWatchstate{
		UserID:   userID,
		MovieID:  movieID,
		Position: input.Position,
		Duration: input.Duration,
		Finished: input.Finished,
	}

	if err := s.watchstateRepo.Upsert(ctx, watchstate); err != nil {
		return nil, fmt.Errorf("save movie watchstate: %w", err)
	}

	return watchstate, nil
}

func (s *UserMovieWatchstateService) GetByMovieID(ctx context.Context, userID int64, movieID string) (*entity.UserMovieWatchstate, error) {
	if userID <= 0 {
		return nil, ErrInvalidMovieWatchstateInput
	}

	movieUUID, err := encoders.DecodeUUID(movieID)
	if err != nil {
		return nil, fmt.Errorf("decode movie id: %w", ErrInvalidMovieWatchstateInput)
	}

	watchstate, err := s.watchstateRepo.GetByUserAndMovieID(ctx, userID, movieUUID)
	if err != nil {
		return nil, fmt.Errorf("get movie watchstate by movie id: %w", err)
	}

	return watchstate, nil
}

func (s *UserMovieWatchstateService) DeleteByMovieID(ctx context.Context, userId int64, movieId string) error {
	movieUUID, err := encoders.DecodeUUID(movieId)
	if err != nil {
		return fmt.Errorf("decode movie id: %w", err)
	}

	err = s.watchstateRepo.DeleteByUserAndMovieID(ctx, userId, movieUUID)
	return err
}

func (s *UserMovieWatchstateService) ListByUserID(ctx context.Context, userID int64) ([]entity.UserMovieWatchstate, error) {
	if userID <= 0 {
		return nil, ErrInvalidMovieWatchstateInput
	}

	watchstates, err := s.watchstateRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list movie watchstates by user id: %w", err)
	}

	return watchstates, nil
}
