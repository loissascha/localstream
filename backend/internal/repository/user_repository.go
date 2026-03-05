package repository

import (
	"context"
	"errors"

	"github.com/loissascha/localstream/internal/entity"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	List(ctx context.Context) ([]entity.User, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
}
