package handler_test

import (
	handler "learning_go/internal/handlers"
	"learning_go/internal/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_LogIn(t *testing.T) {
	testCases := []struct {
		username       string
		password       string
		statusResponse int
	}{
		{
			username:       "no_username",
			password:       "no_password",
			statusResponse: http.StatusNotImplemented,
		},
	}

	for _, val := range testCases {
		t.Run(val.password, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			// Create the original handler
			originalHandler := handler.GetLogs()

			// Wrap the original handler with our logging middleware
			// This demonstrates the middleware pattern where we can add functionality
			// around our existing handler without modifying it
			wrappedHandler := middleware.AuthenticateMiddleware(originalHandler)

			// Use the wrapped handler instead of the original
			wrappedHandler.ServeHTTP(w, r)

			if w.Result().StatusCode != val.statusResponse {
				t.Error("unexpected response")
			}
		})
	}
}
