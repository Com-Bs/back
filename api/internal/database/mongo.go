package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewMongoDB creates a new MongoDB connection
func NewMongoDB(ctx context.Context) (*MongoDB, error) {
	// Get MongoDB URI from environment variable or use default
	uri := os.Getenv("MONGO_URI")
	log.Println("MongoDB URI:", uri)
	if uri == "" {
		// Default connection string for local development
		uri = "mongodb://admin:changeme123@localhost:27017/compis?authSource=admin"
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)
	
	// Set connection timeout
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB successfully")

	// Get database name from environment variable or use default
	dbName := os.Getenv("MONGO_DATABASE")
	if dbName == "" {
		dbName = "compis"
	}

	return &MongoDB{
		Client:   client,
		Database: client.Database(dbName),
	}, nil
}

// Disconnect closes the MongoDB connection
func (m *MongoDB) Disconnect(ctx context.Context) error {
	if err := m.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}
	log.Println("Disconnected from MongoDB")
	return nil
}

// GetCollection returns a collection from the database
func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}