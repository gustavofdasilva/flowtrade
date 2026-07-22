package services

import (
	"account-service/internal/core/domain"
	"account-service/internal/core/ports"
	"fmt"
	"time"

	"github.com/google/uuid"
)

//TODO: Make func 'revokeToken' to also hash before access db or remove hash inside the 'saveToken' to standardize

type AuthService struct {
	userRepo             ports.UserRepository
	authRepo             ports.AuthRepository
	hasher               ports.PasswordHasher
	tokenProvider        ports.TokenProvider
	refreshTokenDuration time.Duration
}

func NewAuthService(userRepo ports.UserRepository, authRepo ports.AuthRepository, hasher ports.PasswordHasher, tokenProvider ports.TokenProvider, refreshTokenDuration time.Duration) *AuthService {
	return &AuthService{
		userRepo:             userRepo,
		authRepo:             authRepo,
		hasher:               hasher,
		tokenProvider:        tokenProvider,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (s *AuthService) Login(email string, password string) (*domain.AuthResult, error) {

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	err = s.hasher.ComparePassword(user.Password, password)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	token, err := s.tokenProvider.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenProvider.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	err = s.saveRefreshToken(user.ID, refreshToken)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return &domain.AuthResult{
		AccessToken:  token,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) Refresh(refreshToken string) (*domain.AuthResult, error) {

	userID, err := s.GetUserIDByRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(*userID)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenProvider.GenerateToken(userID.String(), user.Email)
	if err != nil {
		return nil, err
	}

	hashedToken := s.tokenProvider.HashToken(refreshToken)
	err = s.authRepo.RevokeRefreshToken(hashedToken)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.tokenProvider.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	err = s.saveRefreshToken(user.ID, newRefreshToken)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResult{
		AccessToken:  token,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) Logout(refreshToken string) error {

	hashedToken := s.tokenProvider.HashToken(refreshToken)
	err := s.authRepo.RevokeRefreshToken(hashedToken)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) LogoutAll(userID uuid.UUID) error {
	err := s.authRepo.RevokeAllRefreshTokens(userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GetUserIDByRefreshToken(refreshAccessToken string) (*uuid.UUID, error) {

	hashedToken := s.tokenProvider.HashToken(refreshAccessToken)

	refreshToken, err := s.authRepo.GetRefreshToken(hashedToken)
	if err != nil {
		return nil, fmt.Errorf("error retrieving refresh token: %v", err)
	}

	if !refreshToken.IsValid() {
		return nil, domain.ErrRefreshTokenExpired
	}

	return &refreshToken.ID, nil
}

func (s *AuthService) saveRefreshToken(userID uuid.UUID, refreshToken string) error {

	expiresAt := time.Now().Add(s.refreshTokenDuration)

	hashedToken := s.tokenProvider.HashToken(refreshToken)

	err := s.authRepo.SaveRefreshToken(userID, hashedToken, expiresAt)
	if err != nil {
		return fmt.Errorf("error saving refresh token: %v", err)
	}

	return nil
}
