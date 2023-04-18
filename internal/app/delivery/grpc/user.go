package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	uc_model "github.com/ziyadovea/task_manager/users/internal/app/entity"
	"github.com/ziyadovea/task_manager/users/proto/v1/pb"
)

type userService struct {
	pb.UnimplementedUserServiceServer
	uc UserUsecase
}

func NewUserService(uc UserUsecase) pb.UserServiceServer {
	return userService{
		uc: uc,
	}
}

func (u userService) RegisterUser(ctx context.Context, user *pb.User) (*pb.UserView, error) {
	if user == nil {
		return nil, status.Error(codes.InvalidArgument, "nil user")
	}

	ucUser := ProtoUser2UcUser(user)
	registeredUser, err := u.uc.RegisterUser(ctx, ucUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to register user: %s", err)
	}

	return UcUser2ProtoUserView(registeredUser), nil
}

func (u userService) AuthenticateUser(ctx context.Context, request *pb.AuthenticateUserRequest) (*pb.AuthenticateUserResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "nil request")
	}

	accessToken, refreshToken, err := u.uc.AuthenticateUser(ctx, uc_model.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to auth user: %s", err)
	}

	return &pb.AuthenticateUserResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u userService) RefreshUserToken(ctx context.Context, request *pb.RefreshUserTokenRequest) (*pb.RefreshUserTokenResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "nil request")
	}

	newAccessToken, err := u.uc.RefreshUserToken(ctx, request.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to refresh user token: %s", err)
	}

	return &pb.RefreshUserTokenResponse{
		AccessToken: newAccessToken,
	}, nil
}

func (u userService) ValidateUserToken(ctx context.Context, request *pb.ValidateUserTokenRequest) (*pb.ValidateUserTokenResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "nil request")
	}

	userID, err := u.uc.ValidateUserToken(ctx, request.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to validate user token: %s", err)
	}

	return &pb.ValidateUserTokenResponse{
		UserId: userID,
	}, nil
}

func (u userService) UpdateUser(ctx context.Context, user *pb.User) (*pb.UpdateUserResponse, error) {
	if user == nil {
		return nil, status.Error(codes.InvalidArgument, "nil user")
	}

	ctx = contextWithUserId(ctx)

	updatedCount, err := u.uc.UpdateUser(ctx, ProtoUser2UcUser(user))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to update user: %s", err)
	}

	return &pb.UpdateUserResponse{UpdatedCount: updatedCount}, nil
}

func (u userService) RemoveUser(ctx context.Context, request *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "nil request")
	}

	ctx = contextWithUserId(ctx)

	removedCount, err := u.uc.RemoveUser(ctx, uc_model.User{ID: request.UserId})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to remove user: %s", err)
	}

	return &pb.RemoveUserResponse{RemovedCount: removedCount}, nil
}

func (u userService) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.UserView, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "nil request")
	}

	ctx = contextWithUserId(ctx)

	deletedUser, err := u.uc.GetUser(ctx, uc_model.User{ID: request.UserId})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to get user: %s", err)
	}

	return UcUser2ProtoUserView(deletedUser), nil
}

func (u userService) ListUsers(ctx context.Context, request *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "nil request")
	}

	ctx = contextWithUserId(ctx)

	users, err := u.uc.ListUsers(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to get list users: %s", err)
	}

	pbUserViews := make([]*pb.UserView, len(users))
	for i, u := range users {
		pbUserViews[i] = UcUser2ProtoUserView(u)
	}

	return &pb.ListUsersResponse{
		Users: pbUserViews,
	}, nil
}

func contextWithUserId(ctx context.Context) context.Context {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userIds := md.Get("user_id")
		if len(userIds) > 0 {
			ctx = context.WithValue(ctx, "user_id", userIds[0])
		}
	}
	return ctx
}
