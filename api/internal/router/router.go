package router

import (
	handler "learning_go/internal/handlers"
	"learning_go/internal/middleware"
	"net/http"
)

// Chain applies middlewares to a handler in the correct order
func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func New() http.Handler {
	r := http.NewServeMux()

	// Signup route - allows new users to register
	r.Handle("/signUp", handler.SignUp())

	// Login route - JWT creation handled directly in handler
	r.Handle("/logIn", handler.LogIn())

	// Protected routes that require authentication
	r.Handle("/logs", Chain(
		handler.GetLogs(),
		middleware.AuthenticateMiddleware,    // Verifies JWT token
		middleware.DBLoggingMiddleware("DB"), // Logs the request
	))

	// Compile route with authentication and logging
	r.Handle("/compile", Chain(
		handler.GetFullCompile(),
		middleware.AuthenticateMiddleware,    // Verifies JWT token
		middleware.DBLoggingMiddleware("DB"), // Logs the request
	))

	return r
}
