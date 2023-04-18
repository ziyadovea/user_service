package grpc

import (
	"context"

	uc_model "github.com/ziyadovea/task_manager/users/internal/app/entity"
	"github.com/ziyadovea/task_manager/users/proto/v1/pb"
)

type UserUsecase interface {
	RegisterUser(ctx context.Context, user uc_model.User) (uc_model.User, error)
	AuthenticateUser(ctx context.Context, user uc_model.User) (string, string, error)
	RefreshUserToken(ctx context.Context, refreshToken string) (string, error)
	ValidateUserToken(ctx context.Context, token string) (int64, error)
	UpdateUser(ctx context.Context, user uc_model.User) (int64, error)
	RemoveUser(ctx context.Context, user uc_model.User) (int64, error)
	GetUser(ctx context.Context, user uc_model.User) (uc_model.User, error)
	ListUsers(ctx context.Context) ([]uc_model.User, error)
}

func ProtoUser2UcUser(u *pb.User) uc_model.User {
	return uc_model.User{
		ID:       u.Id,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}

func UcUser2ProtoUserView(u uc_model.User) *pb.UserView {
	return &pb.UserView{
		Id:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
