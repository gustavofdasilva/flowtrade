package handler

import (
	"account-service/internal/core/domain"
	"account-service/internal/core/ports"
	"context"

	pb "github.com/gustavofdasilva/flowtrade/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserGRPCHandler struct {
	pb.UnimplementedAccountServiceServer
	svc ports.UserService
}

func NewUserGRPCHandler(svc ports.UserService) *UserGRPCHandler {
	return &UserGRPCHandler{svc: svc}
}

func (h *UserGRPCHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*emptypb.Empty, error) {
	user := domain.User{
		Username: req.Username,
		Email:    req.Username,
		Password: req.Password,
	}

	_, err := h.svc.Register(user)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}
