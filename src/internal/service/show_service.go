package service

import (
	"context"

	"github.com/loissascha/localstream/internal/encoders"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type ShowService struct {
	showRepo repository.ShowRepository
}

func NewShowService(showRepo repository.ShowRepository) *ShowService {
	return &ShowService{showRepo: showRepo}
}

func (s *ShowService) GetByIDWithMetadata(ctx context.Context, id string) (*repository.ShowSelectItem, error) {
	iduu, err := encoders.DecodeUUID(id)
	if err != nil {
		return nil, err
	}
	return s.showRepo.GetByIDWithMetadata(ctx, iduu)
}

func (s *ShowService) GetByID(ctx context.Context, id string) (*entity.Show, error) {
	iduu, err := encoders.DecodeUUID(id)
	if err != nil {
		return nil, err
	}
	return s.showRepo.GetByID(ctx, iduu)
}

func (s *ShowService) ListLatest(ctx context.Context) ([]repository.ShowSelectItem, error) {
	return s.showRepo.ListLatest(ctx)
}

func (s *ShowService) List(ctx context.Context) ([]repository.ShowSelectItem, error) {
	return s.showRepo.List(ctx)
}

func (s *ShowService) Search(ctx context.Context, query string) ([]repository.ShowSelectItem, error) {
	return s.showRepo.Search(ctx, query)
}
