package model

import (
	"context"
	"errors"
	"learning_go/internal/auth"
	"time"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the database
type User struct {
	// ID is the unique identifier for the user
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	// Username is the unique username for the user
	Username string `json:"username" bson:"username"`
	// Password is the hashed password for the user
	Password string `json:"password" bson:"password"`
	// Email is the email address of the user
	Email string `json:"email" bson:"email"`
	// CreatedAt is the timestamp when the user was created
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	// UpdatedAt is the timestamp when the user was last updated
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// UserService handles user operations with the database
type UserService struct {
	Collection *mongo.Collection
}

// NewUserService creates a new user service
func NewUserService(db *mongo.Database) *UserService {
	return &UserService{
		Collection: db.Collection("users"),
	}
}

// CreateUser stores a new user in the database
func (us *UserService) CreateUser(ctx context.Context, username, email, password string) (*User, error) {
	// validate email, username is validated in getUserByUsername
	if !IsSanitized(email) {
		return nil, errors.New("email contains invalid characters")
	}
	
	// Hash the password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}


	// Check if user already exists
	existingUser, err := us.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}else if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	user := &User{
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert user into database
	result, err := us.Collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	// Set the ID from the insert result
	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

func IsSanitized(s string) bool {
	return !strings.ContainsAny(s, "<>\"'${}[]|\\^`")
}

// GetUserByUsername retrieves a user by username
func (us *UserService) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	// Validate username
	if !IsSanitized(username) {
		return nil, errors.New("username contains invalid characters")
	}
	
	var user User
	filter := bson.M{"username": username}

	err := us.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (us *UserService) GetUserByID(ctx context.Context, id string) (*User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	var user User
	filter := bson.M{"_id": objectID}

	err = us.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetAllUsers retrieves all users from the database
func (us *UserService) GetAllUsers(ctx context.Context) ([]*User, error) {
	cursor, err := us.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*User
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser updates an existing user
func (us *UserService) UpdateUser(ctx context.Context, id string, updates bson.M) (*User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	updates["updated_at"] = time.Now()
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": updates}

	_, err = us.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Return the updated user
	return us.GetUserByID(ctx, id)
}

// DeleteUser deletes a user from the database
func (us *UserService) DeleteUser(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user ID")
	}

	filter := bson.M{"_id": objectID}
	_, err = us.Collection.DeleteOne(ctx, filter)
	return err
}

// ValidateUserPassword validates a user's password
func (us *UserService) ValidateUserPassword(ctx context.Context, username, password string) (bool, error) {
	user, err := us.GetUserByUsername(ctx, username)
	if err != nil {
		return false, err
	}

	// Verify password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, errors.New("invalid password")
	}

	return true, nil
}
