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

type UserHandler struct {
	svc ports.UserService
}

func NewUserHandler(svc ports.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

// TODO: return token in 'data' field
// Register godoc
//
// @Summary Register new user
// @Description Creates a new user account. Requires email, username, and password.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration payload"
// @Success 204 "User created successfully"
// @Failure 400 {object} responses.ErrorResponse "Invalid fields or email/username already in use"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	req := dto.RegisterRequest{}

	c.Bind(&req)

	if req.Email == "" {
		response.Error(c, http.StatusBadRequest, domain.ErrInvalidEmail)
		return
	}
	if req.Password == "" {
		response.Error(c, http.StatusBadRequest, domain.ErrInvalidPassword)
		return
	}
	if req.Username == "" {
		response.Error(c, http.StatusBadRequest, domain.ErrInvalidUsername)
		return
	}

	user := domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	_, err := h.svc.Register(user)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyInUse) {
			response.Error(c, http.StatusBadRequest, err)
			return
		}

		if errors.Is(err, domain.ErrUsernameAlreadyInUse) {
			response.Error(c, http.StatusBadRequest, err)
			return
		}

		slog.Error("Failed to register user", slog.Any("err", err))
		response.Error(c, http.StatusInternalServerError, domain.ErrUnexpectedError)
		return
	}

	c.JSON(http.StatusNoContent, nil)
	return
}

// UpdateUser godoc
//
// @Summary Update current user
// @Description Updates the authenticated user's profile information (email and/or username)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateUserRequest true "Update payload"
// @Success 200 {object} dto.UserResponse "Returns updated user data"
// @Failure 404 {object} responses.ErrorResponse "User not found"
// @Failure 409 {object} responses.ErrorResponse "Email or username already in use"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /users/me [patch]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userIDstr := c.GetString("userID")
	userID := uuid.MustParse(userIDstr)

	req := dto.UpdateUserRequest{}

	c.Bind(&req)

	user := domain.User{
		ID:       userID,
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	updatedUser, err := h.svc.Update(user)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyInUse) {
			response.Error(c, http.StatusConflict, err)
			return
		}

		if errors.Is(err, domain.ErrUsernameAlreadyInUse) {
			response.Error(c, http.StatusConflict, err)
			return
		}

		if errors.Is(err, domain.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, err)
			return
		}

		slog.Error("Failed to update user", slog.Any("err", err))
		response.Error(c, http.StatusInternalServerError, domain.ErrUnexpectedError)
		return
	}

	data := dto.UserResponse{
		ID:       updatedUser.ID,
		Username: updatedUser.Username,
		Email:    updatedUser.Email,
	}

	response.Success(c, http.StatusOK, "User updated successfully", data)
	return
}

// DeleteUser godoc
//
// @Summary Delete current user
// @Description Soft-deletes the authenticated user's account
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 204 "User deleted successfully"
// @Failure 404 {object} responses.ErrorResponse "User not found"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /users/me [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userIDstr := c.GetString("userID")
	userID := uuid.MustParse(userIDstr)

	err := h.svc.Delete(userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, err)
			return
		}

		slog.Error("Failed to delete user", slog.Any("err", err))
		response.Error(c, http.StatusInternalServerError, domain.ErrUnexpectedError)
		return
	}

	c.JSON(http.StatusNoContent, nil)
	return
}

// GetUserInfo godoc
//
// @Summary Get current user info
// @Description Returns the authenticated user's profile information
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse "Returns user data"
// @Failure 404 {object} responses.ErrorResponse "User not found"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /users/me [get]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userIDstr := c.GetString("userID")
	userID := uuid.MustParse(userIDstr)

	user, err := h.svc.GetByID(userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, err)
			return
		}

		slog.Error("Failed to get user info", slog.Any("err", err))
		response.Error(c, http.StatusInternalServerError, domain.ErrUnexpectedError)
		return
	}

	data := dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	response.Success(c, http.StatusOK, "User info retrieved successfully", data)
	return
}
