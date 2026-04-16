package repository

import (
	"context"
	"errors"

	"github.com/loissascha/localstream/internal/entity"
)

var ErrAppSettingsNotFound = errors.New("app settings not found")

type AppSettingsRepository interface {
	Get(ctx context.Context) (*entity.AppSettings, error)
	Update(ctx context.Context, appSettings *entity.AppSettings) error
}
