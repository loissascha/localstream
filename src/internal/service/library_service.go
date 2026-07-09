package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/loissascha/localstream/internal/encoders"
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

func (s *LibraryService) Update(ctx context.Context, id, name, path string, libraryType entity.LibraryType) (*entity.Library, error) {
	trimmedName := strings.TrimSpace(name)
	trimmedPath := strings.TrimSpace(path)
	if trimmedName == "" || trimmedPath == "" {
		return nil, ErrInvalidLibraryInput
	}

	libraryUUID, err := encoders.DecodeUUID(id)
	if err != nil {
		return nil, err
	}

	library, err := s.libraryRepo.GetByID(ctx, libraryUUID)
	if err != nil {
		return nil, err
	}

	library.Name = trimmedName
	library.Path = trimmedPath
	library.LibraryType = libraryType

	err = s.libraryRepo.Update(ctx, library)
	if err != nil {
		return nil, err
	}

	return library, nil
}

func (s *LibraryService) Create(ctx context.Context, name, path, libraryType string) (*entity.Library, error) {
	trimmedName := strings.TrimSpace(name)
	trimmedPath := strings.TrimSpace(path)
	normalizedLibraryType := entity.LibraryType(strings.ToLower(strings.TrimSpace(libraryType)))
	if trimmedName == "" || trimmedPath == "" || !normalizedLibraryType.IsValid() {
		return nil, ErrInvalidLibraryInput
	}

	library := &entity.Library{
		Name:        trimmedName,
		Path:        trimmedPath,
		LibraryType: normalizedLibraryType,
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
