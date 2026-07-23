package handler

import (
	response "api-gateway/internal/adapters/http"
	"api-gateway/internal/adapters/http/dto"
	"api-gateway/internal/core/domain"
	"api-gateway/internal/core/ports"
	"log"
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

func (h *AccountHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	input := domain.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := h.accountClient.Login(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	data := dto.TokensResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		User: dto.UserResponse{
			ID:       res.User.ID,
			Username: res.User.Username,
			Email:    res.User.Email,
		},
	}

	response.Success(c, http.StatusOK, "Login successful", data)
}

func (h *AccountHandler) Refresh(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	input := domain.RefreshTokenInput{
		RefreshToken: req.RefreshToken,
	}

	res, err := h.accountClient.RefreshToken(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	data := dto.TokensResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		User: dto.UserResponse{
			ID:       res.User.ID,
			Username: res.User.Username,
			Email:    res.User.Email,
		},
	}

	response.Success(c, http.StatusOK, "Refresh successful", data)
}

func (h *AccountHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	input := domain.LogoutInput{
		Token: req.Token,
	}

	err := h.accountClient.Logout(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *AccountHandler) LogoutAll(c *gin.Context) {
	userID := c.GetString("userID")

	input := domain.UserIDInput{
		ID: userID,
	}

	err := h.accountClient.LogoutAll(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *AccountHandler) GetMe(c *gin.Context) {
	userID := c.GetString("userID")

	input := domain.UserIDInput{
		ID: userID,
	}

	res, err := h.accountClient.GetMe(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	data := dto.UserResponse{
		ID:       res.ID,
		Username: res.Username,
		Email:    res.Email,
	}

	response.Success(c, http.StatusOK, "User info retrieved", data)
}

func (h *AccountHandler) UpdateUser(c *gin.Context) {
	userID := c.GetString("userID")
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	input := domain.UpdateUserInput{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	log.Println(userID)

	res, err := h.accountClient.UpdateUser(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	data := dto.UserResponse{
		ID:       res.ID,
		Username: res.Username,
		Email:    res.Email,
	}

	response.Success(c, http.StatusOK, "User info updated", data)
}

func (h *AccountHandler) DeleteUser(c *gin.Context) {
	userID := c.GetString("userID")

	input := domain.UserIDInput{
		ID: userID,
	}

	err := h.accountClient.DeleteUser(c.Request.Context(), input)
	if err != nil {
		response.Error(c, grpcStatusToHTTP(err), err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

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
