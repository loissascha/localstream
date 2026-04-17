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

func (s *MovieService) ListLatest(ctx context.Context) ([]repository.MovieSelectItem, error) {
	return s.movieRepo.ListLatest(ctx)
}

func (s *MovieService) List(ctx context.Context) ([]repository.MovieSelectItem, error) {
	return s.movieRepo.List(ctx)
}

func (s *MovieService) GetByIDWithMetadata(ctx context.Context, id string) (*repository.MovieSelectItem, error) {
	uid, err := encoders.DecodeUUID(id)
	if err != nil {
		return nil, err
	}
	return s.movieRepo.GetByIDWithMetadata(ctx, uid)
}

func (s *MovieService) GetById(ctx context.Context, id string) (*entity.Movie, error) {
	uid, err := encoders.DecodeUUID(id)
	if err != nil {
		return nil, err
	}
	return s.movieRepo.GetByID(ctx, uid)
}

func (s *MovieService) Search(ctx context.Context, query string) ([]repository.MovieSelectItem, error) {
	return s.movieRepo.Search(ctx, query)
}
