package model

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	Password       []byte    `json:"-"` // Hide password from JSON
	Email          string    `json:"email"`
	Role           int       `json:"role"`
	TimeCreated    time.Time `json:"timeCreated"`
	TimeLastActive time.Time `json:"timeLastActive"`
}

// UserResources represents resources owned by a user
type UserResources struct {
	UserID   int     `json:"userId"`
	Money    float64 `json:"money"`
	Prestige int     `json:"prestige"`
}

// UserUpdateRequest represents a request to update user data
type UserUpdateRequest struct {
	ID          int    `json:"id"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	OldPassword string `json:"oldPassword,omitempty"`
	NewPassword string `json:"newPassword,omitempty"`
}
