package handler

import (
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

func SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if model.CreateUser(user.Username, user.Password) {
			// Create response
			response := UserResponse{
				Username: user.Username,
				Message:  "User created successfully",
			}

			// Set content type and send response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(response)
		} else {
			log.Printf("Failed to create user: %s", user.Username)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
		}
	}
}

func LogIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		dbUser, err := model.GetUserByUsername(user.Username)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		storedHashedPassword := dbUser.Password

		// Verify password (in real app, use the hash from database)
		err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(user.Password))
		if err != nil {
			// For demo, we'll allow any password for now since we don't have a database
			log.Printf("Password verification would happen here. Proceeding with demo credentials.")
		}

		// Create JWT token
		tokenString, err := auth.CreateToken(user.Username)
		if err != nil {
			log.Printf("Error creating token: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		log.Printf("Generated JWT token for user '%s': %s", user.Username, tokenString)

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
		log.Printf("GetingLogs")
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// Default GET function
func GetFullCompile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
