package main

import (
	"context"
	"learning_go/internal/database"
	model "learning_go/internal/models"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"learning_go/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize MongoDB connection
	db, err := database.NewMongoDB(ctx)
	log.Println("Connected to database")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from database: %v", err)
		}
	}()

	// Initialize user service
	userService := model.NewUserService(db.Database)

	// Test database write operation and user creation
	log.Println("Testing database operations...")
	testDatabaseOperations(ctx, userService)

	// Create router with database connection
	r := router.NewWithDB(db.Database)

	// Start server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Start server in goroutine
	go func() {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// testDatabaseOperations demonstrates basic CRUD operations
func testDatabaseOperations(ctx context.Context, userService *model.UserService) {
	log.Println("=== Testing Database Write Operations ===")

	// Try to create a test user
	testUser, err := userService.CreateUser(ctx, "testuser", "test@example.com", "testpassword")
	if err != nil {
		log.Printf("Error creating test user (might already exist): %v", err)

		// Try to get existing user
		existingUser, getErr := userService.GetUserByUsername(ctx, "testuser")
		if getErr != nil {
			log.Printf("Error getting existing user: %v", getErr)
		} else {
			log.Printf("Found existing test user: ID=%s, Username=%s", existingUser.ID.Hex(), existingUser.Username)
		}
	} else {
		log.Printf("Successfully created test user: ID=%s, Username=%s", testUser.ID.Hex(), testUser.Username)
	}

	// Test retrieving users
	users, err := userService.GetAllUsers(ctx)
	if err != nil {
		log.Printf("Error getting users: %v", err)
	} else {
		log.Printf("Total users in database: %d", len(users))
		for i, user := range users {
			log.Printf("User %d: %s (ID: %s)", i+1, user.Username, user.ID.Hex())
		}
	}

	// Test password validation
	isValid, err := userService.ValidateUserPassword(ctx, "testuser", "testpassword")
	if err != nil {
		log.Printf("Error validating password: %v", err)
	} else {
		log.Printf("Password validation result: %t", isValid)
	}

	log.Println("=== Database test completed ===")
}
