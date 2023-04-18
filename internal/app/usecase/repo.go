package usecase

import (
	"context"

	"github.com/ziyadovea/task_manager/users/internal/app/entity"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user entity.User) (entity.User, error)
	GetUserByID(ctx context.Context, id int64) (entity.User, error)
	GetUserByName(ctx context.Context, name string) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	ListUsers(ctx context.Context) ([]entity.User, error)
	UpdateUserByID(ctx context.Context, user entity.User) (int64, error)
	RemoveUserByID(ctx context.Context, id int64) (int64, error)
}
