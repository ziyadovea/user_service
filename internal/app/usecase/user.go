package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/ziyadovea/task_manager/users/internal/app/entity"
)

// prometheus metric to measure user registration duration
var userRegistrationDuration = promauto.NewSummary(
	prometheus.SummaryOpts{
		Name:       "user_usecase_register_user_duration_seconds_ms",
		Help:       "Summary of User Registration duration",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
)

type userUsecase struct {
	repo          UserRepository
	authenticator Authenticator
}

func NewUserUsecase(repo UserRepository, authenticator Authenticator) userUsecase {
	return userUsecase{
		repo:          repo,
		authenticator: authenticator,
	}
}

func (u userUsecase) RegisterUser(ctx context.Context, user entity.User) (entity.User, error) {
	// measure prometheus metric
	start := time.Now()
	defer func() {
		userRegistrationDuration.Observe(float64(time.Since(start).Milliseconds()))
	}()

	if err := user.Validate(); err != nil {
		return entity.User{}, fmt.Errorf("invalid user: %w", err)
	}

	if err := user.HashPassword(); err != nil {
		return entity.User{}, fmt.Errorf("unable to hash password: %w", err)
	}

	insertedUser, err := u.repo.InsertUser(ctx, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("unable to insert user in repo: %w", err)
	}

	return insertedUser, nil
}

func (u userUsecase) AuthenticateUser(ctx context.Context, user entity.User) (string, string, error) {
	var repoUser entity.User
	if user.Name != "" {
		u, err := u.repo.GetUserByName(ctx, user.Name)
		if err != nil {
			return "", "", fmt.Errorf("unable to get user by name from repo: %w", err)
		}
		repoUser = u
	} else if user.Email != "" {
		u, err := u.repo.GetUserByEmail(ctx, user.Email)
		if err != nil {
			return "", "", fmt.Errorf("unable to get user by email from repo: %w", err)
		}
		repoUser = u
	}

	if err := repoUser.ComparePassword(user.Password); err != nil {
		return "", "", fmt.Errorf("invalid password: %w", err)
	}

	accessToken, err := u.authenticator.CreateAccessToken(repoUser.ID)
	if err != nil {
		return "", "", fmt.Errorf("unable to create access token: %w", err)
	}

	refreshToken, err := u.authenticator.CreateRefreshToken(repoUser.ID)
	if err != nil {
		return "", "", fmt.Errorf("unable to create refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (u userUsecase) RefreshUserToken(ctx context.Context, refreshToken string) (string, error) {
	userID, err := u.authenticator.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("unable to verify refresh token: %w", err)
	}

	accessToken, err := u.authenticator.CreateAccessToken(userID)
	if err != nil {
		return "", fmt.Errorf("unable to create access token: %w", err)
	}

	return accessToken, nil
}

func (u userUsecase) ValidateUserToken(ctx context.Context, token string) (int64, error) {
	return u.authenticator.VerifyAccessToken(token)
}

func (u userUsecase) UpdateUser(ctx context.Context, user entity.User) (int64, error) {
	if err := user.HashPassword(); err != nil {
		return 0, fmt.Errorf("unable to hash password: %w", err)
	}
	return u.repo.UpdateUserByID(ctx, user)
}

func (u userUsecase) RemoveUser(ctx context.Context, user entity.User) (int64, error) {
	return u.repo.RemoveUserByID(ctx, user.ID)
}

func (u userUsecase) GetUser(ctx context.Context, user entity.User) (entity.User, error) {
	if user.ID != 0 {
		repoUser, err := u.repo.GetUserByID(ctx, user.ID)
		if err != nil {
			return entity.User{}, fmt.Errorf("unable to get user by id from repo: %w", err)
		}
		return repoUser, nil
	}

	if user.Name != "" {
		repoUser, err := u.repo.GetUserByName(ctx, user.Name)
		if err != nil {
			return entity.User{}, fmt.Errorf("unable to get user by name from repo: %w", err)
		}
		return repoUser, nil
	}

	if user.Email != "" {
		repoUser, err := u.repo.GetUserByEmail(ctx, user.Email)
		if err != nil {
			return entity.User{}, fmt.Errorf("unable to get user by email from repo: %w", err)
		}
		return repoUser, nil
	}

	return entity.User{}, errors.New("invalid user")
}

func (u userUsecase) ListUsers(ctx context.Context) ([]entity.User, error) {
	return u.repo.ListUsers(ctx)
}
