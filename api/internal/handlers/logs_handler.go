package handler

import (
	"encoding/json"
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

func GetAllLogsByProblemAndUser(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate HTTP method
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract username from context (set by AuthenticateMiddleware)
		username, ok := r.Context().Value("username").(string)
		if !ok || username == "" {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// Extract problem ID from request body
		var requestBody struct {
			ProblemID string `json:"problem_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		problemID := requestBody.ProblemID
		if problemID == "" {
			http.Error(w, "Problem ID is required", http.StatusBadRequest)
			return
		}

		// Initialize log service
		logService := model.NewLogsService(db)

		logs, err := logService.GetLogByProblemID(ctx, problemID, username)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if logs == nil {
			http.Error(w, "No logs found for this problem and user", http.StatusNotFound)
			return
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Encode and send response
	}
}
