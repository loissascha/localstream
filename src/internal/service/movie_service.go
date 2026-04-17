package service

import (
	"context"

	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type MovieService struct {
	movieRepo repository.MovieRepository
}

func NewMovieService(movieRepo repository.MovieRepository) *MovieService {
	return &MovieService{
		movieRepo: movieRepo,
	}
}

func (s *MovieService) ListLatest(ctx context.Context, userID int64) ([]repository.MovieSelectItem, error) {
	return s.movieRepo.ListLatest(ctx, userID)
}

func (s *MovieService) List(ctx context.Context, userID int64) ([]repository.MovieSelectItem, error) {
	return s.movieRepo.List(ctx, userID)
}

func (s *MovieService) GetByIDWithMetadata(ctx context.Context, id string, userID int64) (*repository.MovieSelectItem, error) {
	uid, err := encoders.DecodeUUID(id)
	if err != nil {
		return nil, err
	}
	return s.movieRepo.GetByIDWithMetadata(ctx, uid, userID)
}

func (s *MovieService) GetById(ctx context.Context, id string) (*entity.Movie, error) {
	uid, err := encoders.DecodeUUID(id)
	if err != nil {
		return nil, err
	}
	return s.movieRepo.GetByID(ctx, uid)
}

func (s *MovieService) Search(ctx context.Context, query string, userID int64) ([]repository.MovieSelectItem, error) {
	return s.movieRepo.Search(ctx, query, userID)
}
