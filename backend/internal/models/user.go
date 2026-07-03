package models

import (
	"time"
)

// User represents a user in the system.
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Role         string    `json:"role"`
	Bio          string    `json:"bio"`
	Phone        string    `json:"phone,omitempty"`
	AvatarURL    string    `json:"avatar_url"`
	Languages    []string  `json:"languages"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserCreateRequest is used for creating a new user (registration).
type UserCreateRequest struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role,omitempty"` // defaults to "traveler"
}

// UserUpdateRequest is used for updating an existing user.
type UserUpdateRequest struct {
	Email     string   `json:"email,omitempty"`
	Username  string   `json:"username,omitempty"`
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
	Bio       string   `json:"bio,omitempty"`
	Phone     string   `json:"phone,omitempty"`
	AvatarURL string   `json:"avatar_url,omitempty"`
	Languages []string `json:"languages,omitempty"`
}

// LoginRequest is used for authenticating a user.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse is returned after successful authentication.
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
