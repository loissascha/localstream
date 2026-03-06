package service

import "github.com/loissascha/localstream/internal/repository"

type ShowService struct {
	showRepo repository.ShowRepository
}

func NewShowService(showRepo repository.ShowRepository) *ShowService {
	return &ShowService{showRepo: showRepo}
}
