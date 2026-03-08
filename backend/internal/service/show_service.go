package service

import (
	"context"

	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

type ShowService struct {
	showRepo repository.ShowRepository
}

func NewShowService(showRepo repository.ShowRepository) *ShowService {
	return &ShowService{showRepo: showRepo}
}

func (s *ShowService) List(ctx context.Context) ([]entity.Show, error) {
	return s.showRepo.List(ctx)
}
