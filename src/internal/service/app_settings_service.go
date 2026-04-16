package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/loissascha/localstream/internal/entity"
	"github.com/loissascha/localstream/internal/repository"
)

var ErrInvalidAppSettingsInput = errors.New("invalid app settings input")

type AppSettingsService struct {
	appSettingsRepo repository.AppSettingsRepository
}

func NewAppSettingsService(appSettingsRepo repository.AppSettingsRepository) *AppSettingsService {
	return &AppSettingsService{appSettingsRepo: appSettingsRepo}
}

func (s *AppSettingsService) Get(ctx context.Context) (*entity.AppSettings, error) {
	appSettings, err := s.appSettingsRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get app settings: %w", err)
	}

	return appSettings, nil
}

func (s *AppSettingsService) Update(ctx context.Context, executeLibraryWatcher bool, libraryWatcherIntervalSeconds int) (*entity.AppSettings, error) {
	if libraryWatcherIntervalSeconds <= 0 {
		return nil, ErrInvalidAppSettingsInput
	}

	appSettings, err := s.appSettingsRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get app settings: %w", err)
	}

	appSettings.ExecuteLibraryWatcher = executeLibraryWatcher
	appSettings.LibraryWatcherIntervalSeconds = libraryWatcherIntervalSeconds

	if err := s.appSettingsRepo.Update(ctx, appSettings); err != nil {
		return nil, fmt.Errorf("update app settings: %w", err)
	}

	return appSettings, nil
}
