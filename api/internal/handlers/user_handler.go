package handler

import (
	"context"
	"encoding/json"
	"learning_go/internal/auth"
	"learning_go/internal/cache"
	model "learning_go/internal/models"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UserResponse represents the response after successful signup
type UserResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

var ctx = context.Background()

var compileCache = cache.NewCompileCache(24 * time.Hour)

func SignUp(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate HTTP method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse JSON request body
		var user model.User

		log.Printf("Received request to create user: %s", r.Body)

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		// Validate input
		if user.Username == "" || user.Password == "" || user.Email == "" {
			http.Error(w, "Username and password are required", http.StatusBadRequest)
			return
		}

		userService := model.NewUserService(db)
		createdUser, err := userService.CreateUser(ctx, user.Username, user.Email, user.Password)
		if err != nil {
			if err.Error() == "user already exists" {
				log.Printf("User already exists: %s", user.Username)
				http.Error(w, "User already exists", http.StatusConflict)
				return
			} else if err.Error() == "username contains invalid characters" {
				log.Printf("Invalid characters in username: %s", user.Username)
				http.Error(w, "Username contains invalid characters", http.StatusBadRequest)
				return
			} else if err.Error() == "email contains invalid characters" {
				log.Printf("Invalid characters in email: %s", user.Username)
				http.Error(w, "email contains invalid characters", http.StatusBadRequest)
				return
			}

			log.Printf("Failed to create user: %v", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		// Create response
		token, err := auth.CreateToken(createdUser.Username)
		if err != nil {
			log.Printf("Error creating token: %v", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		response := UserResponse{
			Message: "User created successfully",
			Token:   token,
		}

		// Set content type and send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func LogIn(db *mongo.Database) http.HandlerFunc {
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

		userService := model.NewUserService(db)
		dbUser, err := userService.GetUserByUsername(ctx, user.Username)
		if err != nil {
			if err.Error() == "user not found" {
				http.Error(w, "User not found", http.StatusUnauthorized)
			} else {
				http.Error(w, "Invalid characters in username", http.StatusBadRequest)
			}
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
			"message": "Login successful",
			"token":   tokenString,
		}

		// Set content type and send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
