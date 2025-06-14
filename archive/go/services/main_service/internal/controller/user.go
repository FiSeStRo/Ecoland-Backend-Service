package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/model"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/service"
	"github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service/internal/utils"
)

// UserController handles HTTP requests related to users
type UserController struct {
	userService service.UserService
	authService service.AuthService
}

// NewUserController creates a new user controller
func NewUserController(userService service.UserService, authService service.AuthService) *UserController {
	return &UserController{
		userService: userService,
		authService: authService,
	}
}

// RegisterRoutes registers all user-related routes
func (c *UserController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /user/resources", c.GetUserResources)
	mux.HandleFunc("GET /user/info", c.GetUserInfo)
	mux.HandleFunc("PATCH /user/update", c.UpdateUserInfo)
}

// GetUserResources handles retrieving user resources
func (c *UserController) GetUserResources(w http.ResponseWriter, r *http.Request) {
	claims, err := c.authService.ValidateAuthentication(r)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	resources, err := c.userService.GetUserResources(claims.UserId)
	if err != nil {
		http.Error(w, "Failed to retrieve user resources", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, resources, http.StatusOK)
}

// GetUserInfo handles retrieving user profile information
func (c *UserController) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	claims, err := c.authService.ValidateAuthentication(r)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		userIDStr = strconv.Itoa(claims.UserId)
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err = c.userService.ValidateUser(claims.UserId, userID); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := c.userService.GetUserInfo(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve user information", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, user, http.StatusOK)
}

// UpdateUserInfo handles updating user profile information
func (c *UserController) UpdateUserInfo(w http.ResponseWriter, req *http.Request) {
	claims, err := c.authService.ValidateAuthentication(req)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var body model.UserUpdateRequest
	if err := json.NewDecoder(req.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if body.ID != claims.UserId {
		http.Error(w, "Invalid requestBody", http.StatusBadRequest)
	}

	if body.Email != "" && !utils.IsEmail(body.Email) {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	user, err := c.userService.UpdateUserInfo(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, user, http.StatusOK)
}
