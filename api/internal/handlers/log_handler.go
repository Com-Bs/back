package handler

import (
	"encoding/json"
	"learning_go/internal/middleware"
	model "learning_go/internal/models"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

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
			"success":        true,
			"solutions":      solutions,
			"totalSolutions": len(solutions),
		}

		// Set content type and send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
