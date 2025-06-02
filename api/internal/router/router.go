package router

import (
	handler "learning_go/internal/handlers"
	"learning_go/internal/middleware"
	"net/http"
)

func New() *http.ServeMux {
	r := http.NewServeMux()

	// and from right to left after the request
	r.Handle("/api/compile", Chain(
		// First: Authentication (pre-request check)
		middleware.AuthenticateMiddleware,
		// Second: Logging (pre-request setup and post-request logging)
		middleware.DBLoggingMiddleware("DB"),
	)(handler.GetLogs()))

	return r
}

type Middleware func(http.Handler) http.Handler

func Chain(middlewares ...Middleware) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		return handler
	}
}
