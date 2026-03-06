package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/loissascha/localstream/internal/repository"
)

type ShowService struct {
	showRepo repository.ShowRepository
}

func NewShowService(showRepo repository.ShowRepository) *ShowService {
	return &ShowService{showRepo: showRepo}
}

func (s *ShowService) ListForLibrary(ctx context.Context, libraryID uuid.UUID) {
}
