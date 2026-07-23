package handler

import (
	"account-service/internal/adapters/http/dto"
	"account-service/internal/adapters/http/response"
	"account-service/internal/core/domain"
	"account-service/internal/core/ports"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	svc ports.AuthService
}

func NewAuthHandler(svc ports.AuthService) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}

// Login godoc
//
// @Summary Login user
// @Description Authenticates a user and returns access and refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login payload"
// @Success 200 {object} dto.TokensResponse "Returns access_token, refresh_token, and user info"
// @Failure 400 {object} responses.ErrorResponse "Missing email or password"
// @Failure 401 {object} responses.ErrorResponse "Invalid credentials"
// @Failure 404 {object} responses.ErrorResponse "User not found"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	req := dto.LoginRequest{}

	c.Bind(&req)

	if req.Email == "" {
		response.Error(c, http.StatusBadRequest, domain.ErrInvalidEmail)
		return
	}
	if req.Password == "" {
		response.Error(c, http.StatusBadRequest, domain.ErrInvalidPassword)
		return
	}

	refreshToken, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, err)
			return
		}

		if errors.Is(err, domain.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, err)
			return
		}

		slog.Error("Failed to login user", slog.Any("err", err))
		response.Error(c, http.StatusInternalServerError, domain.ErrUnexpectedError)
		return
	}

	data := dto.TokensResponse{
		AccessToken:  refreshToken.AccessToken,
		RefreshToken: refreshToken.RefreshToken,
		User: dto.UserResponse{
			ID:       refreshToken.User.ID,
			Username: refreshToken.User.Username,
			Email:    refreshToken.User.Email,
		},
	}

	response.Success(c, http.StatusOK, "Login successful", data)
	return
}

// RefreshToken godoc
//
// @Summary Refresh access token
// @Description Rotates the access and refresh token pair using a valid refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token payload"
// @Success 200 {object} dto.TokensResponse "Returns new access_token, refresh_token, and user info"
// @Failure 400 {object} responses.ErrorResponse "Missing refresh token"
// @Failure 401 {object} responses.ErrorResponse "Invalid, expired, or not found refresh token"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	req := dto.RefreshTokenRequest{}

	c.Bind(&req)

	if req.RefreshToken == "" {
		response.Error(c, http.StatusBadRequest, domain.ErrInvalidRefreshToken)
		return
	}

	refreshToken, err := h.svc.Refresh(req.RefreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrRefreshTokenNotFound) || errors.Is(err, domain.ErrInvalidRefreshToken) || errors.Is(err, domain.ErrRefreshTokenExpired) {
			response.Error(c, http.StatusUnauthorized, err)
			return
		}

		slog.Error("Failed to refresh token", slog.Any("err", err))
		response.Error(c, http.StatusInternalServerError, domain.ErrUnexpectedError)
		return
	}

	data := dto.TokensResponse{
		AccessToken:  refreshToken.AccessToken,
		RefreshToken: refreshToken.RefreshToken,
		User: dto.UserResponse{
			ID:       refreshToken.User.ID,
			Username: refreshToken.User.Username,
			Email:    refreshToken.User.Email,
		},
	}

	response.Success(c, http.StatusOK, "Refresh successful", data)
	return
}

// Logout godoc
//
// @Summary Logout user
// @Description Invalidates a specific refresh token, ending the current session
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.RefreshTokenRequest true "Refresh token to invalidate"
// @Success 200 {object} responses.SuccessResponse "Logout successful"
// @Failure 400 {object} responses.ErrorResponse "Missing refresh token"
// @Failure 401 {object} responses.ErrorResponse "Invalid, expired, or not found refresh token"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	req := dto.RefreshTokenRequest{}

	c.Bind(&req)

	if req.RefreshToken == "" {
		response.Error(c, http.StatusBadRequest, domain.ErrInvalidRefreshToken)
		return
	}

	err := h.svc.Logout(req.RefreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrRefreshTokenNotFound) || errors.Is(err, domain.ErrInvalidRefreshToken) || errors.Is(err, domain.ErrRefreshTokenExpired) {
			response.Error(c, http.StatusUnauthorized, err)
			return
		}

		slog.Error("Failed to logout user", slog.Any("err", err))
		response.Error(c, http.StatusInternalServerError, domain.ErrUnexpectedError)
		return
	}

	response.Success(c, http.StatusOK, "Logout successful", nil)
	return
}

// LogoutAll godoc
//
// @Summary Logout from all devices
// @Description Invalidates all refresh tokens for the authenticated user
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.SuccessResponse "Logout successful"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized or invalid session"
// @Router /auth/logout-all [post]
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	userIDstr := c.GetString("userID")
	userID := uuid.MustParse(userIDstr)

	err := h.svc.LogoutAll(userID)
	if err != nil {
		if errors.Is(err, domain.ErrRefreshTokenNotFound) || errors.Is(err, domain.ErrInvalidRefreshToken) || errors.Is(err, domain.ErrRefreshTokenExpired) {
			response.Error(c, http.StatusUnauthorized, err)
			return
		}

		slog.Error("Failed to logout all sessions", slog.Any("err", err))
		response.Error(c, http.StatusInternalServerError, domain.ErrUnexpectedError)
		return
	}

	response.Success(c, http.StatusOK, "Logout successful", nil)
	return
}
