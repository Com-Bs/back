package handler

import (
	"context"
	"encoding/json"
	"learning_go/internal/auth"
	"learning_go/internal/cache"
	"learning_go/internal/middleware"
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
			"message": "Login successful",
			"token":   tokenString,
		}

		// Set content type and send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// Default GET function
func GetLogs(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate HTTP method
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		logService := model.NewLogsService(db)
		logs, err := logService.GetAllLogs(ctx)

		if err != nil {
			log.Printf("Failed to retrieve logs: %v", err)
			http.Error(w, "Failed to retrieve logs", http.StatusInternalServerError)
			return
		}
		if len(logs) == 0 {
			http.Error(w, "No logs found", http.StatusNotFound)
			return
		}

		// Set content type and send response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)

		w.WriteHeader(http.StatusOK)
	}
}

// GetUserSolutions returns user's previous solutions for a problem
func GetUserSolutions(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate HTTP method
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract problem ID from URL path
		problemID := r.PathValue("id")
		if problemID == "" {
			http.Error(w, "Problem ID is required", http.StatusBadRequest)
			return
		}

		// Get username from context (set by auth middleware)
		username, ok := r.Context().Value(middleware.UsernameKey).(string)
		if !ok {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		// Query user solutions from logs
		logsService := model.NewLogsService(db)
		solutions, err := logsService.GetUserSolutionsByProblem(ctx, username, problemID)
		if err != nil {
			http.Error(w, "Failed to retrieve solutions", http.StatusInternalServerError)
			return
		}

		// Create response
		response := map[string]interface{}{
			"success":       true,
			"solutions":     solutions,
			"totalSolutions": len(solutions),
		}

		// Set content type and send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
