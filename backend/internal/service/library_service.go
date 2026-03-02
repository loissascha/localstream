package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

var ErrInvalidLibraryInput = errors.New("invalid library input")

type LibraryService struct {
	libraryRepo repository.LibraryRepository
}

func NewLibraryService(libraryRepo repository.LibraryRepository) *LibraryService {
	return &LibraryService{libraryRepo: libraryRepo}
}

func (s *LibraryService) Create(ctx context.Context, name, path string) (*entity.Library, error) {
	trimmedName := strings.TrimSpace(name)
	trimmedPath := strings.TrimSpace(path)
	if trimmedName == "" || trimmedPath == "" {
		return nil, ErrInvalidLibraryInput
	}

	library := &entity.Library{
		Name: trimmedName,
		Path: trimmedPath,
	}

	if err := s.libraryRepo.Create(ctx, library); err != nil {
		return nil, fmt.Errorf("create library: %w", err)
	}

	return library, nil
}

func (s *LibraryService) List(ctx context.Context) ([]entity.Library, error) {
	libraries, err := s.libraryRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list libraries: %w", err)
	}

	return libraries, nil
}
