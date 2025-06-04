package handler

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
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
	Username string `json:"username"`
	Message  string `json:"message"`
}

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
		userService := model.NewUserService(db)
		createdUser, err := userService.CreateUser(ctx, user.Username, user.Email, user.Password)
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

		ctx := context.Background()
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
			"token":   tokenString,
			"message": "Login successful",
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

		ctx := context.Background()
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

// Default GET function
func CreateLogs(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate HTTP method
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

	}
}

// GetFullCompile handles code compilation requests with caching
func GetFullCompile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate HTTP method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Create hash of request body for caching
		hash := sha256.Sum256(body)
		hashStr := hex.EncodeToString(hash[:])

		// Check cache first
		if cached, exists := compileCache.Get(hashStr); exists {
			log.Printf("Cache hit for compile request: %s", hashStr)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(cached.StatusCode)
			w.Write(cached.ResponseBody)
			return
		}

		// Create HTTP client
		client := &http.Client{}

		// Create request to compile service with the body
		req, err := http.NewRequest("POST", "http://172.16.30.3:3001/runCompile", bytes.NewBuffer(body))
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Failed to send request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response", http.StatusInternalServerError)
			return
		}

		// Cache successful responses
		if resp.StatusCode == http.StatusOK {
			compileCache.Set(hashStr, respBody, resp.StatusCode)
			log.Printf("Cached compile response for request: %s", hashStr)
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)

		// Write response back to client
		w.Write(respBody)
	}
}
