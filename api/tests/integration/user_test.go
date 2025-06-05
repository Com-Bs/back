package integration

import (
	"context"
	"learning_go/internal/database"
	"learning_go/internal/router"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

var testDB *mongo.Database

// discardLogger implements io.Writer and discards all writes
type discardLogger struct{}

func (d *discardLogger) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// testLogger is a custom logger that only logs test-related messages
type testLogger struct {
	*testing.T
}

func (tl *testLogger) Printf(format string, v ...interface{}) {
	tl.T.Logf(format, v...)
}

func TestMain(m *testing.M) {
	// Set up test database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Discard all standard logs during tests
	log.SetOutput(&discardLogger{})
	log.SetFlags(0)

	mongoDB, err := database.NewMongoDB(ctx)
	if err != nil {
		panic(err)
	}
	testDB = mongoDB.Database

	// Run tests
	code := m.Run()

	// Clean up
	if err := mongoDB.Disconnect(ctx); err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestLogInRoute(t *testing.T) {
	// Create test logger
	logger := &testLogger{t}
	handler := router.NewWithDB(testDB)

	// Step 2: iterate over each TestCase in unit.logIn
	for _, tc := range LogIn {
		t.Run(tc.Name, func(t *testing.T) {
			logger.Printf("Running test: %s", tc.Name)
			// Build the HTTP request
			var req *http.Request
			if tc.Body != "" {
				req = httptest.NewRequest(tc.Method, tc.URL, strings.NewReader(tc.Body))
			} else {
				req = httptest.NewRequest(tc.Method, tc.URL, nil)
			}

			// Set headers (e.g. Content-Type: application/json)
			for k, v := range tc.Headers {
				req.Header.Set(k, v)
			}

			// Record the response
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tc.ExpectedStatus {
				t.Errorf(
					"Test %q: expected status %d, got %d. Body=%q",
					tc.Name, tc.ExpectedStatus, rr.Code, rr.Body.String(),
				)
			}

			// If ExpectedBody is non-empty, verify it appears
			if tc.ExpectedBody != "" {
				body := rr.Body.String()
				if !strings.Contains(body, tc.ExpectedBody) {
					t.Errorf(
						"Test %q: expected body to contain %q, but got %q",
						tc.Name, tc.ExpectedBody, body,
					)
				}
			}
		})
	}
}

func TestSignUpRoute(t *testing.T) {
	// Create test logger
	logger := &testLogger{t}
	handler := router.NewWithDB(testDB)

	for _, tc := range SignUp {
		t.Run(tc.Name, func(t *testing.T) {
			logger.Printf("Running test: %s", tc.Name)
			var req *http.Request
			if tc.Body != "" {
				req = httptest.NewRequest(tc.Method, tc.URL, strings.NewReader(tc.Body))
			} else {
				req = httptest.NewRequest(tc.Method, tc.URL, nil)
			}
			for k, v := range tc.Headers {
				req.Header.Set(k, v)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.ExpectedStatus {
				t.Errorf(
					"Test %q: expected status %d, got %d. Body=%q",
					tc.Name, tc.ExpectedStatus, rr.Code, rr.Body.String(),
				)
			}
			if tc.ExpectedBody != "" {
				body := rr.Body.String()
				if !strings.Contains(body, tc.ExpectedBody) {
					t.Errorf(
						"Test %q: expected body to contain %q, but got %q",
						tc.Name, tc.ExpectedBody, body,
					)
				}
			}
		})
	}
}

func TestLogsRoute(t *testing.T) {
	// Create test logger
	logger := &testLogger{t}
	handler := router.NewWithDB(testDB)

	for _, tc := range GetLogs {
		t.Run(tc.Name, func(t *testing.T) {
			logger.Printf("Running test: %s", tc.Name)
			var req *http.Request
			if tc.Body != "" {
				req = httptest.NewRequest(tc.Method, tc.URL, strings.NewReader(tc.Body))
			} else {
				req = httptest.NewRequest(tc.Method, tc.URL, nil)
			}
			for k, v := range tc.Headers {
				req.Header.Set(k, v)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.ExpectedStatus {
				t.Errorf(
					"Test %q: expected status %d, got %d. Body=%q",
					tc.Name, tc.ExpectedStatus, rr.Code, rr.Body.String(),
				)
			}
			if tc.ExpectedBody != "" {
				body := rr.Body.String()
				if !strings.Contains(body, tc.ExpectedBody) {
					t.Errorf(
						"Test %q: expected body to contain %q, but got %q",
						tc.Name, tc.ExpectedBody, body,
					)
				}
			}
		})
	}
}

func TestGetProblems(t *testing.T) {
	// Create test logger
	logger := &testLogger{t}
	handler := router.NewWithDB(testDB)

	for _, tc := range GetProblems {
		t.Run(tc.Name, func(t *testing.T) {
			logger.Printf("Running test: %s", tc.Name)
			var req *http.Request
			if tc.Body != "" {
				req = httptest.NewRequest(tc.Method, tc.URL, strings.NewReader(tc.Body))
			} else {
				req = httptest.NewRequest(tc.Method, tc.URL, nil)
			}
			for k, v := range tc.Headers {
				req.Header.Set(k, v)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.ExpectedStatus {
				t.Errorf(
					"Test %q: expected status %d, got %d. Body=%q",
					tc.Name, tc.ExpectedStatus, rr.Code, rr.Body.String(),
				)
			}
			if tc.ExpectedBody != "" {
				body := rr.Body.String()
				if !strings.Contains(body, tc.ExpectedBody) {
					t.Errorf(
						"Test %q: expected body to contain %q, but got %q",
						tc.Name, tc.ExpectedBody, body,
					)
				}
			}
		})
	}
}

func TestGetProblemByID(t *testing.T) {
	// Create test logger
	logger := &testLogger{t}
	handler := router.NewWithDB(testDB)

	for _, tc := range GetProblemByID {
		t.Run(tc.Name, func(t *testing.T) {
			logger.Printf("Running test: %s", tc.Name)
			var req *http.Request
			if tc.Body != "" {
				req = httptest.NewRequest(tc.Method, tc.URL, strings.NewReader(tc.Body))
			} else {
				req = httptest.NewRequest(tc.Method, tc.URL, nil)
			}
			for k, v := range tc.Headers {
				req.Header.Set(k, v)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.ExpectedStatus {
				t.Errorf(
					"Test %q: expected status %d, got %d. Body=%q",
					tc.Name, tc.ExpectedStatus, rr.Code, rr.Body.String(),
				)
			}
			if tc.ExpectedBody != "" {
				body := rr.Body.String()
				if !strings.Contains(body, tc.ExpectedBody) {
					t.Errorf(
						"Test %q: expected body to contain %q, but got %q",
						tc.Name, tc.ExpectedBody, body,
					)
				}
			}
		})
	}
}
