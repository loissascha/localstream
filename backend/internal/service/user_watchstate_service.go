package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

var ErrInvalidWatchstateInput = errors.New("invalid watchstate input")

type SaveWatchstateInput struct {
	ShowID    string
	SeasonID  string
	EpisodeID string
	Position  float64
	Duration  float64
	Finished  bool
}

type UserWatchstateService struct {
	watchstateRepo repository.UserWatchstateRepository
}

func NewUserWatchstateService(watchstateRepo repository.UserWatchstateRepository) *UserWatchstateService {
	return &UserWatchstateService{watchstateRepo: watchstateRepo}
}

func (s *UserWatchstateService) Save(ctx context.Context, userId int64, input SaveWatchstateInput) (*entity.UserWatchstate, error) {
	if userId <= 0 {
		return nil, ErrInvalidWatchstateInput
	}
	if input.Position < 0 || input.Duration < 0 {
		return nil, ErrInvalidWatchstateInput
	}

	showID, err := encoders.DecodeUUID(input.ShowID)
	if err != nil {
		return nil, fmt.Errorf("decode show id: %w", ErrInvalidWatchstateInput)
	}

	seasonID, err := encoders.DecodeUUID(input.SeasonID)
	if err != nil {
		return nil, fmt.Errorf("decode season id: %w", ErrInvalidWatchstateInput)
	}

	episodeID, err := encoders.DecodeUUID(input.EpisodeID)
	if err != nil {
		return nil, fmt.Errorf("decode episode id: %w", ErrInvalidWatchstateInput)
	}

	watchstate := &entity.UserWatchstate{
		UserID:    userId,
		ShowID:    showID,
		SeasonID:  seasonID,
		EpisodeID: episodeID,
		Position:  input.Position,
		Duration:  input.Duration,
		Finished:  input.Finished,
	}

	if err := s.watchstateRepo.Upsert(ctx, watchstate); err != nil {
		return nil, fmt.Errorf("save watchstate: %w", err)
	}

	return watchstate, nil
}

func (s *UserWatchstateService) GetByEpisodeID(ctx context.Context, userId int64, episodeId string) (*entity.UserWatchstate, error) {
	if userId <= 0 {
		return nil, ErrInvalidWatchstateInput
	}

	episodeUUID, err := encoders.DecodeUUID(episodeId)
	if err != nil {
		return nil, fmt.Errorf("decode episode id: %w", ErrInvalidWatchstateInput)
	}

	watchstate, err := s.watchstateRepo.GetByUserAndEpisodeID(ctx, userId, episodeUUID)
	if err != nil {
		return nil, fmt.Errorf("get watchstate by episode id: %w", err)
	}

	return watchstate, nil
}

func (s *UserWatchstateService) ListLatestByShow(ctx context.Context, userId int64) ([]entity.UserWatchstate, error) {
	if userId <= 0 {
		return nil, ErrInvalidWatchstateInput
	}

	watchstates, err := s.watchstateRepo.ListLatestByShowForUserID(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("list latest watchstates by show: %w", err)
	}

	return watchstates, nil
}
