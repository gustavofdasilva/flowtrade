package ports

import (
	"api-gateway/internal/core/domain"
	"context"
)

type AccountClient interface {
	Register(ctx context.Context, input domain.RegisterInput) error
	Login(ctx context.Context, input domain.LoginInput) (*domain.AuthOutput, error)
	RefreshToken(ctx context.Context, input domain.RefreshTokenInput) (*domain.AuthOutput, error)
	Logout(ctx context.Context, input domain.LogoutInput) error
	LogoutAll(ctx context.Context, input domain.UserIDInput) error
	UpdateUser(ctx context.Context, input domain.UpdateUserInput) (*domain.UserOutput, error)
	DeleteUser(ctx context.Context, input domain.UserIDInput) error
	GetMe(ctx context.Context, input domain.UserIDInput) (*domain.UserOutput, error)
}
