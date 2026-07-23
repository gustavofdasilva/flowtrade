package handler

import (
	"api-gateway/internal/adapters/http/dto"
	"api-gateway/internal/adapters/http/response"
	"api-gateway/internal/core/domain"
	"api-gateway/internal/core/ports"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountHandler struct {
	accountClient ports.AccountClient
}

func NewAccountHandler(client ports.AccountClient) *AccountHandler {
	return &AccountHandler{accountClient: client}
}

func (h *AccountHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	input := domain.RegisterInput{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	err := h.accountClient.Register(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// func (h *AccountHandler) Login(c *gin.Context) {
// 	var req dto.LoginRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		response.Error(c, http.StatusBadRequest, err)
// 		return
// 	}

// 	input := domain.LoginInput{
// 		Email:    req.Email,
// 		Password: req.Password,
// 	}

// 	res, err := h.accountClient.Login(c.Request.Context(), input)
// 	if err != nil {
// 		response.Error(c, grpcStatusToHTTP(err), err)
// 		return
// 	}

// 	data := dto.TokensResponse{
// 		AccessToken:  res.AccessToken,
// 		RefreshToken: res.RefreshToken,
// 		User: dto.UserResponse{
// 			ID:       res.User.ID,
// 			Username: res.User.Username,
// 			Email:    res.User.Email,
// 		},
// 	}

// 	response.Success(c, http.StatusOK, "Login successful", data)
// }

func grpcStatusToHTTP(err error) int {
	switch status.Code(err) {
	case codes.NotFound:
		return 404
	case codes.InvalidArgument:
		return 400
	case codes.Unauthenticated:
		return 401
	case codes.PermissionDenied:
		return 403
	case codes.AlreadyExists:
		return 409
	case codes.ResourceExhausted:
		return 429
	default:
		return 500
	}
}
