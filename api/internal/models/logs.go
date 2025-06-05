package model

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Logs struct {
	// ID is the unique identifier for the log entry
	ID primitive.ObjectID `bson:"_id,omitempty"`
	// UserID is the ID of the user associated with the log entry
	UserID *string `bson:"user_id,omitempty"`
	// Method is the HTTP method used in the request (e.g., GET, POST)
	Method string `bson:"method"`
	// Path is the URL path of the request
	Path string `bson:"path"`
	// Query is the query parameters of the request
	ResponseStatus int `bson:"status"`
	// ResponseStatus is the HTTP status code returned in the response
	Duration time.Duration `bson:"duration"`
	// Duration is the time taken to process the request
	Body string `bson:"body"`
	// Body is the request body, typically a JSON string
	Problem primitive.ObjectID `bson:"case,omitempty"`
	// Problem is the ID of the problem associated with the log entry
	IP string `bson:"ip"`
	// IP is the IP address of the user making the request
	CreatedAt time.Time `bson:"created_at"`
	// CreatedAt is the timestamp when the log entry was created
}

type LogsService struct {
	Collection *mongo.Collection
}

func NewLogsService(db *mongo.Database) *LogsService {
	return &LogsService{
		Collection: db.Collection("logs"),
	}
}

func (ls *LogsService) CreateLog(ctx context.Context, logs *Logs) error {

	log.Print("Creating log entry...")

	as, err := ls.Collection.InsertOne(context.TODO(), logs)

	log.Print("Log entry created successfully", as)

	if err != nil {
		return err
	}
	return nil
}

func (ls *LogsService) GetLogsByUserID(ctx context.Context, userID string) ([]*Logs, error) {
	filter := primitive.M{"user_id": userID}
	cursor, err := ls.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logsList []*Logs
	for cursor.Next(ctx) {
		var log Logs
		if err := cursor.Decode(&log); err != nil {
			return nil, err
		}
		logsList = append(logsList, &log)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return logsList, nil
}

func (ls *LogsService) GetAllLogs(ctx context.Context) ([]*Logs, error) {
	cursor, err := ls.Collection.Find(ctx, primitive.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logsList []*Logs
	for cursor.Next(ctx) {
		var log Logs
		if err := cursor.Decode(&log); err != nil {
			return nil, err
		}
		logsList = append(logsList, &log)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return logsList, nil
}

func (ls *LogsService) GetLogByHash(ctx context.Context, hash string) (*Logs, error) {
	var log Logs
	filter := primitive.M{"body": hash}
	err := ls.Collection.FindOne(ctx, filter).Decode(&log)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}

func (ls *LogsService) GetLogByProblemID(ctx context.Context, id string, userID string) (*Logs, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var log Logs
	filter := primitive.M{
		"problem": objectID,
		"user":    userID,
	}

	err = ls.Collection.FindOne(ctx, filter).Decode(&log)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &log, nil
}
