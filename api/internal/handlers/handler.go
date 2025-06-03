package handler

import (
	"context"
	"encoding/json"
	"learning_go/internal/auth"
	model "learning_go/internal/models"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// UserResponse represents the response after successful signup
type UserResponse struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

// UserHandler holds the user service
type UserHandler struct {
	UserService *model.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *model.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (uh *UserHandler) SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate HTTP method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse JSON request body
		var user model.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		// Validate input
		if user.Username == "" || user.Password == "" {
			http.Error(w, "Username and password are required", http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		createdUser, err := uh.UserService.CreateUser(ctx, user.Username, user.Email, user.Password)
		if err != nil {
			log.Printf("Failed to create user: %v", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		// Create response
		response := UserResponse{
			Username: createdUser.Username,
			Message:  "User created successfully",
		}

		// Set content type and send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func (uh *UserHandler) LogIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate HTTP method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse JSON request body
		var user model.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		// Validate input
		if user.Username == "" || user.Password == "" {
			http.Error(w, "Username and password are required", http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		dbUser, err := uh.UserService.GetUserByUsername(ctx, user.Username)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		storedHashedPassword := dbUser.Password

		err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(user.Password))
		if err != nil {
			log.Printf("Password verification failed")
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Create JWT token
		tokenString, err := auth.CreateToken(user.Username)
		if err != nil {
			log.Printf("Error creating token: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Create response with token
		response := map[string]string{
			"token":   tokenString,
			"message": "Login successful",
		}

		// Set content type and send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func LogOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// Default GET function
func GetLogs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate HTTP method
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		log.Printf("GetingLogs")
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// Default GET function
func GetFullCompile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate HTTP method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusNotImplemented)
	}
}
