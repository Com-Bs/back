package model

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	ResponseBody string `bson:"response_body,omitempty"`
	// ResponseBody is the response body from the compile service
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

// UserSolution represents a user's solution attempt from logs
type UserSolution struct {
	ID           string    `json:"id"`
	ProblemID    string    `json:"problemId"`
	UserID       string    `json:"userId"`
	Code         string    `json:"code"`
	Status       string    `json:"status"`
	SubmittedAt  string    `json:"submittedAt"`
	ExecutionTime string  `json:"executionTime"`
}

// GetUserSolutionsByProblem retrieves user's compile attempts for a specific problem
func (ls *LogsService) GetUserSolutionsByProblem(ctx context.Context, userID, problemID string) ([]*UserSolution, error) {
	problemObjectID, err := primitive.ObjectIDFromHex(problemID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"user_id": userID,
		"path":    "/compile",
		"case":    problemObjectID,
	}

	// Sort by creation time (newest first)
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})
	
	cursor, err := ls.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var solutions []*UserSolution
	for cursor.Next(ctx) {
		var logEntry Logs
		if err := cursor.Decode(&logEntry); err != nil {
			continue // Skip invalid entries
		}

		// Convert log entry to UserSolution
		solution := &UserSolution{
			ID:        logEntry.ID.Hex(),
			ProblemID: problemID,
			UserID:    userID,
			Code:      logEntry.Body,
			SubmittedAt: logEntry.CreatedAt.Format(time.RFC3339),
			ExecutionTime: logEntry.Duration.String(),
		}

		// Determine status based on response body
		if logEntry.ResponseBody != "" {
			// Parse the compile response to determine actual status
			var compileResponse struct {
				Result []struct {
					Status         string `json:"status"`
					Output         []int  `json:"output"`
					ExpectedOutput []int  `json:"expectedOutput"`
				} `json:"result"`
				Status string `json:"status"`
				Error  string `json:"error"`
				Line   int    `json:"line"`
				Column int    `json:"column"`
			}
			
			if err := json.Unmarshal([]byte(logEntry.ResponseBody), &compileResponse); err == nil {
				// Check if there's a compilation error (syntax error, etc.)
				if compileResponse.Error != "" {
					solution.Status = "failed"
				} else {
					// Count passed and failed test cases
					passedCount := 0
					totalCount := len(compileResponse.Result)
					
					for _, result := range compileResponse.Result {
						if result.Status == "Success" {
							passedCount++
						}
					}
					
					if passedCount == 0 {
						solution.Status = "failed"
					} else if passedCount == totalCount {
						solution.Status = "passed"
					} else {
						solution.Status = "partial"
					}
				}
			} else {
				// Default to failed if response parsing fails
				solution.Status = "failed"
			}
		} else {
			// Default to failed if no response body
			solution.Status = "failed"
		}

		solutions = append(solutions, solution)
	}

	return solutions, nil
}
