package model

import (
	"errors"
	"learning_go/internal/auth"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the database
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // This will store the hashed password
}

// CreateUser stores a new user in the database
func CreateUser(username, password string) bool {
	// Hash the password
	_, err := auth.HashPassword(password)

	return err == nil
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(username string) (User, error) {
	// Simulate a database lookup
	user := User{
		ID:       1,
		Username: username,
		Password: "$2a$10$EIX/5Z1z5Q8e1b7f9j3u6O0k5F4y5Z1z5Q8e1b7f9j3u6O0k5F4y5Z", // Example hashed password
	}

	return user, nil
}

// ValidateUserPassword validates a user's password
func ValidateUserPassword(username, password string) (bool, error) {
	user, err := GetUserByUsername(username)

	if err != nil {
		return false, err
	}

	// Verify password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, errors.New("invalid password")
	}

	return true, nil
}
