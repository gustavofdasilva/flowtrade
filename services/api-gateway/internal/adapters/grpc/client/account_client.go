// api-gateway/adapters/grpc/client/account_client.go
package grpcclient

import (
	"api-gateway/internal/core/domain"
	"context"

	pb "github.com/gustavofdasilva/flowtrade/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AccountClient struct {
	client pb.AccountServiceClient
}

func NewAccountClient(addr string) *AccountClient {
	conn, _ := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return &AccountClient{client: pb.NewAccountServiceClient(conn)}
}

func (c *AccountClient) Register(ctx context.Context, input domain.RegisterInput) error {
	_, err := c.client.Register(ctx, &pb.RegisterRequest{
		Email:    input.Email,
		Password: input.Password,
		Username: input.Username,
	})

	return err
}

func (c *AccountClient) Login(ctx context.Context, input domain.LoginInput) (*domain.AuthOutput, error) {
	token, err := c.client.Login(ctx, &pb.LoginRequest{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}

	return &domain.AuthOutput{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User: domain.UserOutput{
			ID:       token.User.Id,
			Username: token.User.Username,
			Email:    token.User.Email,
		},
	}, nil
}

func (c *AccountClient) RefreshToken(ctx context.Context, input domain.RefreshTokenInput) (*domain.AuthOutput, error) {
	token, err := c.client.RefreshToken(ctx, &pb.RefreshTokenRequest{
		RefreshToken: input.RefreshToken,
	})
	if err != nil {
		return nil, err
	}

	return &domain.AuthOutput{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User: domain.UserOutput{
			ID:       token.User.Id,
			Username: token.User.Username,
			Email:    token.User.Email,
		},
	}, nil
}

func (c *AccountClient) Logout(ctx context.Context, input domain.LogoutInput) error {
	//TODO: Remake pb funcs to support LogoutRequest
	_, err := c.client.Logout(ctx, &pb.RefreshTokenRequest{
		RefreshToken: input.Token,
	})

	return err
}

func (c *AccountClient) LogoutAll(ctx context.Context, input domain.UserIDInput) error {
	_, err := c.client.LogoutAll(ctx, &pb.LogoutAllRequest{
		UserId: input.ID,
	})

	return err
}

func (c *AccountClient) UpdateUser(ctx context.Context, input domain.UpdateUserInput) (*domain.UserOutput, error) {
	user, err := c.client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:       input.ID,
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}

	newUser := domain.UserOutput{
		ID:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}

	return &newUser, err
}

func (c *AccountClient) GetMe(ctx context.Context, input domain.UserIDInput) (*domain.UserOutput, error) {
	userInfo, err := c.client.GetMe(ctx, &pb.UserIDRequest{
		Id: input.ID,
	})
	if err != nil {
		return nil, err
	}

	user := domain.UserOutput{
		ID:       userInfo.Id,
		Username: userInfo.Username,
		Email:    userInfo.Email,
	}

	return &user, err
}

func (c *AccountClient) DeleteUser(ctx context.Context, input domain.UserIDInput) error {
	_, err := c.client.DeleteUser(ctx, &pb.UserIDRequest{
		Id: input.ID,
	})

	return err
}
