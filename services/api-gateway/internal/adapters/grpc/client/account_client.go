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
