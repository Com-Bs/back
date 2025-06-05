package model

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Logs struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	UserID         *string            `bson:"user_id,omitempty"`
	Method         string             `bson:"method"`
	Path           string             `bson:"path"`
	ResponseStatus int                `bson:"status"`
	Duration       time.Duration      `bson:"duration"`
	Body           string             `bson:"body"`
	Problem        primitive.ObjectID `bson:"case,omitempty"` // Optional field for the case
	Success        bool               `bson:"success"`
	IP             string             `bson:"ip"`
	CreatedAt      time.Time          `bson:"created_at"`
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
