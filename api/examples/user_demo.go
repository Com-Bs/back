package main

import (
	"context"
	"fmt"
	"learning_go/internal/database"
	model "learning_go/internal/models"
	"log"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	// Load .env file
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Create context
	ctx := context.Background()

	// Connect to MongoDB
	db, err := database.NewMongoDB(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from database: %v", err)
		}
	}()

	// Create user service
	userService := model.NewUserService(db.Database)

	// Demo: Create a new user
	fmt.Println("=== Creating a new user ===")
	newUser, err := userService.CreateUser(ctx, "johndoe", "john@example.com", "password123")
	if err != nil {
		log.Printf("Error creating user: %v", err)
	} else {
		fmt.Printf("Created user: %+v\n", newUser)
	}

	// Demo: Create another user
	fmt.Println("\n=== Creating another user ===")
	newUser2, err := userService.CreateUser(ctx, "janedoe", "jane@example.com", "securepass456")
	if err != nil {
		log.Printf("Error creating user: %v", err)
	} else {
		fmt.Printf("Created user: %+v\n", newUser2)
	}

	// Demo: Get user by username
	fmt.Println("\n=== Getting user by username ===")
	user, err := userService.GetUserByUsername(ctx, "johndoe")
	if err != nil {
		log.Printf("Error getting user: %v", err)
	} else {
		fmt.Printf("Found user: %+v\n", user)
	}

	// Demo: Get all users
	fmt.Println("\n=== Getting all users ===")
	users, err := userService.GetAllUsers(ctx)
	if err != nil {
		log.Printf("Error getting users: %v", err)
	} else {
		fmt.Printf("Total users: %d\n", len(users))
		for i, u := range users {
			fmt.Printf("User %d: ID=%s, Username=%s, Email=%s\n",
				i+1, u.ID.Hex(), u.Username, u.Email)
		}
	}

	// Demo: Validate password
	fmt.Println("\n=== Validating password ===")
	isValid, err := userService.ValidateUserPassword(ctx, "johndoe", "password123")
	if err != nil {
		log.Printf("Error validating password: %v", err)
	} else {
		fmt.Printf("Password valid: %t\n", isValid)
	}

	// Demo: Update user
	if len(users) > 0 {
		fmt.Println("\n=== Updating user ===")
		userID := users[0].ID.Hex()
		updates := bson.M{
			"email": "newemail@example.com",
		}
		updatedUser, err := userService.UpdateUser(ctx, userID, updates)
		if err != nil {
			log.Printf("Error updating user: %v", err)
		} else {
			fmt.Printf("Updated user: %+v\n", updatedUser)
		}
	}

	// Demo: Get user by ID
	if len(users) > 0 {
		fmt.Println("\n=== Getting user by ID ===")
		userID := users[0].ID.Hex()
		userByID, err := userService.GetUserByID(ctx, userID)
		if err != nil {
			log.Printf("Error getting user by ID: %v", err)
		} else {
			fmt.Printf("Found user by ID: %+v\n", userByID)
		}
	}

	// Demo: Test connection by inserting a test document
	fmt.Println("\n=== Testing database write operation ===")
	testCollection := db.Database.Collection("test")
	testDoc := bson.M{
		"message":   "Hello MongoDB!",
		"timestamp": time.Now(),
		"test":      true,
	}

	result, err := testCollection.InsertOne(ctx, testDoc)
	if err != nil {
		log.Printf("Error inserting test document: %v", err)
	} else {
		fmt.Printf("Successfully inserted test document with ID: %v\n", result.InsertedID)

		// Clean up test document
		_, err = testCollection.DeleteOne(ctx, bson.M{"_id": result.InsertedID})
		if err != nil {
			log.Printf("Error deleting test document: %v", err)
		} else {
			fmt.Println("Test document cleaned up successfully")
		}
	}

	fmt.Println("\n=== Demo completed ===")
}
