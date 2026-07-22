package services

import (
	"account-service/internal/core/domain"
	"account-service/internal/core/ports"

	"github.com/google/uuid"
)

type UserService struct {
	repo   ports.UserRepository
	hasher ports.PasswordHasher
}

func NewUserService(repo ports.UserRepository, hasher ports.PasswordHasher) *UserService {
	return &UserService{
		repo:   repo,
		hasher: hasher,
	}
}

func (s *UserService) Register(user domain.User) (*domain.User, error) {
	exists, err := s.repo.IsEmailAlreadyInUse(user.Email, nil)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, domain.ErrEmailAlreadyInUse
	}

	exists, err = s.repo.IsUsernameAlreadyInUse(user.Username, nil)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, domain.ErrUsernameAlreadyInUse
	}

	hashedPassword, err := s.hasher.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword

	newUser, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *UserService) Update(user domain.User) (*domain.User, error) {

	exists, err := s.repo.IsEmailAlreadyInUse(user.Email, &user.ID)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, domain.ErrEmailAlreadyInUse
	}

	exists, err = s.repo.IsUsernameAlreadyInUse(user.Username, &user.ID)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, domain.ErrUsernameAlreadyInUse
	}

	hashedPassword, err := s.hasher.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword

	err = s.repo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return &user, nil
}

func (s *UserService) Delete(id uuid.UUID) error {

	err := s.repo.DeleteUserByID(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetByID(id uuid.UUID) (*domain.User, error) {

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return user, nil
}
