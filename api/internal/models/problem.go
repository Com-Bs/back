package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TestCase struct {
	Input  string `json:"input" bson:"input"`
	Output string `json:"output" bson:"output"`
}

type Problem struct {
	// A problem has a title, a description (contains examples), a difficulty, hints and several test cases (each test case has an input and an output)
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title string             `json:"title" bson:"title"`
	// Description is a markdown string that contains examples
	Description string `json:"description" bson:"description"`
	// Difficulty is an enum that represents the difficulty of the problem
	Difficulty string `json:"difficulty" bson:"difficulty"`
	// TestCases is a list of test cases
	TestCases []TestCase `json:"test_cases" bson:"test_cases"`
	// Function name for the problem
	FunctionName string `json:"function_name" bson:"function_name"`
	// Arguments/parameters for the function
	Arguments []ParamType `json:"arguments" bson:"arguments"`
	// CreatedAt is the date and time the problem was created
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	// UpdatedAt is the date and time the problem was last updated
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type ParamType struct {
	Name string `json:"name" bson:"name"`
	Type string `json:"type" bson:"type"` // e.g., "int", "string", "float"
}

type ProblemService struct {
	Collection *mongo.Collection
}

func NewProblemService(db *mongo.Database) *ProblemService {
	return &ProblemService{Collection: db.Collection("problems")}
}

// Usar un singleton para crear las instancias en la base de datos. Creamos un servicio, y al inicializar el api inicializamos el servicio.

// Definimos CRUD para problemas

// GetProblem retrieves a problem by its ID from the database
func (ps *ProblemService) GetProblemByID(ctx context.Context, id string) (*Problem, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var problem Problem
	err = ps.Collection.FindOne(ctx, primitive.M{"_id": objectID}).Decode(&problem)
	if err != nil {
		return nil, err
	}
	return &problem, nil
}

func (ps *ProblemService) GetAllProblems(ctx context.Context) ([]*Problem, error) {
	cursor, err := ps.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var problems []*Problem
	if err = cursor.All(ctx, &problems); err != nil {
		return nil, err
	}
	return problems, nil
}
