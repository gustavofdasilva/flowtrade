package handler

import (
	"account-service/internal/core/domain"
	"account-service/internal/core/ports"
	"context"

	"github.com/google/uuid"
	pb "github.com/gustavofdasilva/flowtrade/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AccountGRPCHandler struct {
	pb.UnimplementedAccountServiceServer
	authSvc ports.AuthService
	userSvc ports.UserService
}

func NewAccountGRPCHandler(authSvc ports.AuthService, userSvc ports.UserService) *AccountGRPCHandler {
	return &AccountGRPCHandler{
		authSvc: authSvc,
		userSvc: userSvc,
	}
}

func (h *AccountGRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	token, err := h.authSvc.Login(req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	user := pb.UserInfo{
		Id:       token.User.ID.String(),
		Username: token.User.Username,
		Email:    token.User.Email,
	}

	return &pb.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User:         &user,
	}, nil
}

func (h *AccountGRPCHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.LoginResponse, error) {

	token, err := h.authSvc.Refresh(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	user := pb.UserInfo{
		Id:       token.User.ID.String(),
		Username: token.User.Username,
		Email:    token.User.Email,
	}

	return &pb.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User:         &user,
	}, nil
}

func (h *AccountGRPCHandler) Logout(ctx context.Context, req *pb.RefreshTokenRequest) (*emptypb.Empty, error) {

	err := h.authSvc.Logout(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (h *AccountGRPCHandler) LogoutAll(ctx context.Context, req *pb.LogoutAllRequest) (*emptypb.Empty, error) {

	err := h.authSvc.LogoutAll(uuid.MustParse(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (h *AccountGRPCHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*emptypb.Empty, error) {
	user := domain.User{
		Username: req.Username,
		Email:    req.Username,
		Password: req.Password,
	}

	_, err := h.userSvc.Register(user)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (h *AccountGRPCHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserInfo, error) {
	user := domain.User{
		ID:       uuid.MustParse(req.Id),
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	newUser, err := h.userSvc.Update(user)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &pb.UserInfo{
		Id:       user.ID.String(),
		Username: newUser.Username,
		Email:    newUser.Email,
	}, nil
}

func (h *AccountGRPCHandler) GetMe(ctx context.Context, req *pb.UserIDRequest) (*pb.UserInfo, error) {
	user, err := h.userSvc.GetByID(uuid.MustParse(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &pb.UserInfo{
		Id:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (h *AccountGRPCHandler) DeleteUser(ctx context.Context, req *pb.UserIDRequest) (*emptypb.Empty, error) {
	err := h.userSvc.Delete(uuid.MustParse(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}
