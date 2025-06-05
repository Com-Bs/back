package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"learning_go/internal/auth"
	model "learning_go/internal/models"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Custom response writer that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Context key for username
type contextKey string

const (
	usernameKey contextKey = "username"
	bodyKey     contextKey = "body"
	fullBody    contextKey = "fullBody"
)

// BodyCaptureMiddleware captures the request body and makes it available in the context
func BodyCaptureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only capture body for login and signup
		if r.URL.Path == "/logIn" || r.URL.Path == "/signUp" {
			// Read the body
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				// Restore the body for the next handler
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				// Store the body in context
				ctx := context.WithValue(r.Context(), bodyKey, bodyBytes)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}

// DBLoggingMiddleware returns a middleware that logs HTTP requests to a database
func DBLoggingMiddleware(db *mongo.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Pre-request operations
			start := time.Now()

			// Create a custom response writer to capture the status code
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Execute the next handler (the actual request)
			next.ServeHTTP(rw, r)

			// Get username from context or request body
			var username string
			if ctxUsername, ok := r.Context().Value(usernameKey).(string); ok {
				username = ctxUsername
			} else if r.URL.Path == "/logIn" || r.URL.Path == "/signUp" {
				// For login/signup, try to get username from captured body
				if bodyBytes, ok := r.Context().Value(bodyKey).([]byte); ok {
					var loginData struct {
						Username string `json:"username"`
					}
					if err := json.Unmarshal(bodyBytes, &loginData); err == nil {
						username = loginData.Username
					}
				}
			}

			// Post-request operations
			logEntry := model.Logs{
				UserID:         &username,
				Method:         r.Method,
				Path:           r.URL.Path,
				ResponseStatus: rw.statusCode,
				Duration:       time.Since(start),
				IP:             r.RemoteAddr,
				CreatedAt:      time.Now(),
			}

			if r.URL.Path != "/logIn" && r.URL.Path != "/signUp" {
				// Get code from original request body
				if bodyBytes, ok := r.Context().Value(fullBody).([]byte); ok {
					var body struct {
						Code string `json:"code"`
					}
					if err := json.Unmarshal(bodyBytes, &body); err == nil {
						logEntry.Body = body.Code
					}
					logEntry.ResponseStatus = rw.statusCode
				}
			}

			ctx := context.Background()
			logService := model.NewLogsService(db)

			err := logService.CreateLog(ctx, &logEntry)

			if err != nil {
				log.Printf("Failed to log request: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			log.Printf("Request completed: %v", logEntry)
		})
	}
}

// AuthenticateMiddleware returns a middleware that handles authentication
func AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// For login and signup endpoints, we don't need to verify the token
		if r.URL.Path == "/logIn" || r.URL.Path == "/signUp" {
			next.ServeHTTP(w, r)
			return
		}

		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Extract the token from the Authorization header
		// Format: "Bearer <token>"
		tokenString := authHeader[7:] // Remove "Bearer " prefix

		// Verify the token and get username
		username, err := auth.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add username to request context
		ctx := context.WithValue(r.Context(), usernameKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RepeatedRequestMiddleware(db *mongo.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Read and hash the request body
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {

				// store the body in context for later use
				ctx := context.WithValue(r.Context(), fullBody, bodyBytes)
				r = r.WithContext(ctx)

				// Restore the body for the next handler
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// Hash the body
				hash, err := auth.HashPassword(string(bodyBytes))
				if err == nil {
					// Check if this request was already made
					logService := model.NewLogsService(db)
					if existingLog, err := logService.GetLogByHash(r.Context(), hash); err == nil && existingLog != nil {
						// Request was already made, return the previous response
						w.WriteHeader(existingLog.ResponseStatus)

						return
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
