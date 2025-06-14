package service

import (
	"errors"
	"fmt"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService defines the interface for user-related operations
type UserService interface {
	GetUserInfo(userID int) (*model.User, error)
	GetUserResources(userID int) (*model.UserResources, error)
	UpdateUserInfo(req model.UserUpdateRequest) (*model.User, error)
	ValidateUser(userID, requestedID int) error
}

// userService implements the UserService interface
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetUserInfo retrieves user information
func (s *userService) GetUserInfo(userID int) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	user.Password = nil

	return user, nil
}

// GetUserResources retrieves user resources
func (s *userService) GetUserResources(userID int) (*model.UserResources, error) {
	resources, err := s.userRepo.GetUserResources(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user resources: %w", err)
	}

	return resources, nil
}

// UpdateUserInfo updates user information
func (s *userService) UpdateUserInfo(req model.UserUpdateRequest) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if req.NewPassword != "" {
		if err = bcrypt.CompareHashAndPassword(user.Password, []byte(req.OldPassword)); err != nil {
			return nil, errors.New("incorrect password")
		}
		newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.MinCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = newHash
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	user.Password = nil

	return user, nil
}

func (s *userService) ValidateUser(userID, requestedID int) error {
	// Currently only allow users to access their own data
	// In the future, this could include role-based checks
	if userID != requestedID {
		return errors.New("unauthorized access to user data")
	}

	return nil
}
