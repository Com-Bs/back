package handler

import (
	"context"
	"encoding/json"
	model "learning_go/internal/models"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type ProblemResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Description string `json:"description"`
	Difficulty string `json:"difficulty"`
	Hints     []string `json:"hints"`
	TestCases []TestCase `json:"test_cases"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TestCase struct {
	Input string `json:"input"`
	Output string `json:"output"`
}

func GetProblemByID(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate HTTP method
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract ID from URL path parameter
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "Problem ID is required", http.StatusBadRequest)
			return
		}

		// Initialize problem service
		problemService := model.NewProblemService(db)
		
		// Get problem from database
		problem, err := problemService.GetProblemByID(context.Background(), id)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "Problem not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Encode and send response
		if err := json.NewEncoder(w).Encode(problem); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func GetAllProblems(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate HTTP method
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Initialize problem service
		problemService := model.NewProblemService(db)
		
		// Get all problems from database
		problems, err := problemService.GetAllProblems(context.Background())
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Encode and send response
		if err := json.NewEncoder(w).Encode(problems); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}