package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type EpisodeService struct {
	episodeRepo repository.EpisodeRepository
	seasonRepo  repository.SeasonRepository
}

func NewEpisodeService(episodeRepo repository.EpisodeRepository, seasonRepo repository.SeasonRepository) *EpisodeService {
	return &EpisodeService{episodeRepo: episodeRepo, seasonRepo: seasonRepo}
}

func (s *EpisodeService) ListBySeasonID(ctx context.Context, seasonId string) ([]entity.Episode, error) {
	seasonUUID, err := encoders.DecodeUUID(seasonId)
	if err != nil {
		return nil, fmt.Errorf("parse season id: %w", err)
	}

	episodes, err := s.episodeRepo.ListBySeasonID(ctx, seasonUUID)
	if err != nil {
		return nil, fmt.Errorf("list episodes by season id: %w", err)
	}

	return episodes, nil
}

func (s *EpisodeService) GetByID(ctx context.Context, episodeId string) (*repository.EpisodeWithMetadata, error) {
	id, err := encoders.DecodeUUID(episodeId)
	if err != nil {
		return nil, err
	}

	return s.episodeRepo.GetByID(ctx, id)
}

func (s *EpisodeService) getShowIdForSeasonId(ctx context.Context, seasonId uuid.UUID) (uuid.UUID, error) {
	season, err := s.seasonRepo.GetByID(ctx, seasonId)
	if err != nil {
		return uuid.Nil, err
	}
	showId := season.ShowID
	return showId, nil
}

func (s *EpisodeService) DeleteByID(ctx context.Context, episodeId string) error {
	id, err := encoders.DecodeUUID(episodeId)
	if err != nil {
		return err
	}

	err = s.episodeRepo.DeleteByID(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *EpisodeService) GetNextEpisode(ctx context.Context, episodeId string) (*entity.Episode, error) {
	id, err := encoders.DecodeUUID(episodeId)
	if err != nil {
		return nil, err
	}

	episode, err := s.episodeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	number := episode.Number
	seasonId := episode.SeasonID
	showId, err := s.getShowIdForSeasonId(ctx, seasonId)
	if err != nil {
		return nil, err
	}

	episodes, err := s.episodeRepo.ListBySeasonID(ctx, seasonId)
	if err != nil {
		return nil, err
	}

	// if there is a next episode within the current season -> return it
	useNext := false
	for _, e := range episodes {
		if e.Number == number {
			useNext = true
			continue
		}
		if useNext {
			return &e, nil
		}
	}

	// if there is a next season -> take the first episode of it
	nextSeasonId, err := s.getNextSeasonId(ctx, seasonId, showId)
	if err != nil {
		return nil, err
	}
	if nextSeasonId == uuid.Nil {
		return nil, nil
	}
	episodes, err = s.episodeRepo.ListBySeasonID(ctx, nextSeasonId)
	if err != nil {
		return nil, err
	}
	if len(episodes) > 0 {
		return &episodes[0], nil
	}

	return nil, nil
}

func (s *EpisodeService) getNextSeasonId(ctx context.Context, seasonId, showId uuid.UUID) (uuid.UUID, error) {
	seasons, err := s.seasonRepo.ListByShowID(ctx, showId)
	if err != nil {
		return uuid.Nil, err
	}

	useNext := false
	for _, se := range seasons {
		if se.ID == seasonId {
			useNext = true
			continue
		}
		if useNext {
			return se.ID, nil
		}
	}
	return uuid.Nil, nil
}
