package mariadb

import (
	"database/sql"
	"fmt"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/repository"
)

// userRepository implements the UserRepository interface
type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// CreateUser adds a new user to the database
func (r *userRepository) CreateUser(user *model.User) (int, error) {
	query := `INSERT INTO users(username, password, email) VALUES (?,?,?)`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("prepare statement failed: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Username, user.Password, user.Email)
	if err != nil {
		return 0, fmt.Errorf("insert user failed: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert ID failed: %w", err)
	}

	return int(id), nil
}

// CreateUserResources initializes resources for a new user
func (r *userRepository) CreateUserResources(userID int, money float64, prestige int) error {
	query := `INSERT INTO user_resources(user_id, money, prestige) VALUES(?,?,?)`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare statement failed: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, money, prestige)
	if err != nil {
		return fmt.Errorf("insert user resources failed: %w", err)
	}

	return nil
}

// CheckUserExists checks if a username is already taken
func (r *userRepository) CheckUserExists(username string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE username = ?`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("prepare statement failed: %w", err)
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("query failed: %w", err)
	}

	return count > 0, nil
}

// CheckEmailExists checks if an email is already registered
func (r *userRepository) CheckEmailExists(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ?`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("prepare statement failed: %w", err)
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("query failed: %w", err)
	}

	return count > 0, nil
}

// GetUserByUsername retrieves user by username
func (r *userRepository) GetUserByUsername(username string) (*model.User, error) {
	query := `SELECT id, username, password, email FROM users WHERE username = ?`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare statement failed: %w", err)
	}
	defer stmt.Close()

	var user model.User
	err = stmt.QueryRow(username).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &user, nil
}

// GetUserByID retrieves user by ID
func (r *userRepository) GetUserByID(id int) (*model.User, error) {
	query := `SELECT id, username, password, email FROM users WHERE id = ?`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare statement failed: %w", err)
	}
	defer stmt.Close()

	var user model.User
	err = stmt.QueryRow(id).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &user, nil
}

// GetUserResources retrieves resources for a specific user
func (r *userRepository) GetUserResources(userID int) (*model.UserResources, error) {
	query := `SELECT money, prestige FROM user_resources WHERE user_id = ?`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare statement failed: %w", err)
	}
	defer stmt.Close()

	var resources model.UserResources
	resources.UserID = userID

	err = stmt.QueryRow(userID).Scan(&resources.Money, &resources.Prestige)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user resources not found: %w", err)
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &resources, nil
}

// UpdateUser updates user information
func (r *userRepository) UpdateUser(user *model.User) error {
	query := `UPDATE users SET username = ?, password = ?, email = ? WHERE id = ?`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare statement failed: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Password, user.Email, user.ID)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	return nil
}

// UpdateUserResources updates a user's resources
func (r *userRepository) UpdateUserResources(resources *model.UserResources) error {
	query := `UPDATE user_resources SET money = ?, prestige = ? WHERE user_id = ?`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare statement failed: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(resources.Money, resources.Prestige, resources.UserID)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	return nil
}
