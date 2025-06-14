package repository

import "github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/model"

// UserRepository defines the interface for user-related data access
type UserRepository interface {
	CreateUser(user *model.User) (int, error)
	CreateUserResources(userID int, money float64, prestige int) error
	CheckUserExists(username string) (bool, error)
	CheckEmailExists(email string) (bool, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(id int) (*model.User, error)
	GetUserResources(userID int) (*model.UserResources, error)
	UpdateUser(user *model.User) error
	UpdateUserResources(resources *model.UserResources) error
}
