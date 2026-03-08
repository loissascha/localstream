package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type ShowService struct {
	showRepo repository.ShowRepository
}

func NewShowService(showRepo repository.ShowRepository) *ShowService {
	return &ShowService{showRepo: showRepo}
}

func (s *ShowService) GetByID(ctx context.Context, id string) (*entity.Show, error) {
	iduu, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.showRepo.GetByID(ctx, iduu)
}

func (s *ShowService) List(ctx context.Context) ([]entity.Show, error) {
	return s.showRepo.List(ctx)
}
