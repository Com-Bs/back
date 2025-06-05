package router

import (
	handler "learning_go/internal/handlers"
	"learning_go/internal/middleware"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

// Chain applies middlewares to a handler in the correct order
func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func NewWithDB(db *mongo.Database) http.Handler {
	r := http.NewServeMux()

	// Signup route - POST method for user registration
	r.Handle("POST /signUp", Chain(
		handler.SignUp(db),
		middleware.BodyCaptureMiddleware,
		middleware.DBLoggingMiddleware(db),
	))

	// Login route - POST method for authentication
	r.Handle("POST /logIn", Chain(
		handler.LogIn(db),
		middleware.BodyCaptureMiddleware,
		middleware.DBLoggingMiddleware(db),
	))

	// Protected routes that require authentication
	// GET method for retrieving logs
	r.Handle("GET /logs", Chain(
		handler.GetLogs(db),
		middleware.AuthenticateMiddleware,  // Verifies JWT token
		middleware.DBLoggingMiddleware(db), // Logs the request
	))

	// POST method for code compilation
	r.Handle("POST /compile", Chain(
		handler.GetFullCompile(db),
		middleware.AuthenticateMiddleware,  // Verifies JWT token
		middleware.DBLoggingMiddleware(db), // Logs the request
	))

	// Problem routes
	// GET method for retrieving all problems
	r.Handle("GET /problems", Chain(
		handler.GetAllProblems(db),
		middleware.AuthenticateMiddleware, // Verifies JWT token
	))

	// GET method for retrieving a specific problem by ID
	r.Handle("GET /problems/{id}", Chain(
		handler.GetProblemByID(db),
		middleware.AuthenticateMiddleware, // Verifies JWT token
	))

	r.Handle("Get /problems/{id}/historic", Chain(
		handler.GetAllLogsByProblemAndUser(db),
		middleware.AuthenticateMiddleware, // Verifies JWT token
	))
	return r
}
