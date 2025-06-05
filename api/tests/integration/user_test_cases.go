package integration

import (
	"fmt"
	"math/rand/v2"
)

// TestCase holds one request + expected response.
type TestCase struct {
	Name           string
	Method         string
	URL            string
	Body           string            // if empty, no request body is sent
	Headers        map[string]string // e.g. "Authorization": "Bearer <token>"
	ExpectedStatus int
	// ExpectedBody can be a full exact match or just a substring to look for.
	// If empty, only status code is checked.
	ExpectedBody string
}

var LogIn = []TestCase{
	{
		Name:           "LogIn with valid credentials",
		Method:         "POST",
		URL:            "/logIn",
		Body:           `{"username": "testuser", "password": "testpassword"}`,
		Headers:        map[string]string{"Content-Type": "application/json"},
		ExpectedStatus: 200,
		ExpectedBody:   `{"message":"Login successful","token":`,
	},
	{
		Name:           "LogIn with invalid credentials",
		Method:         "POST",
		URL:            "/logIn",
		Body:           `{"username": "testuser", "password": "wrongpassword"}`,
		Headers:        map[string]string{"Content-Type": "application/json"},
		ExpectedStatus: 401,
		ExpectedBody:   "Invalid credentials",
	},
	{
		Name:           "LogIn with user not found",
		Method:         "POST",
		URL:            "/logIn",
		Body:           `{"username": "nonexistentuser", "password": "testpassword"}`,
		Headers:        map[string]string{"Content-Type": "application/json"},
		ExpectedStatus: 401,
		ExpectedBody:   "User not found",
	},
}

var randomTest = rand.IntN(1000000)

var SignUp = []TestCase{
	{
		Name:           "SignUp with valid data",
		Method:         "POST",
		URL:            "/signUp",
		Body:           fmt.Sprintf(`{"username": "test%dusername", "password": "testpassword", "email": "test%d@email.com"}`, randomTest, randomTest),
		Headers:        map[string]string{"Content-Type": "application/json"},
		ExpectedStatus: 201,
		ExpectedBody:   `{"message":"User created successfully","token":`,
	},
	{
		Name:           "SignUp with missing fields",
		Method:         "POST",
		URL:            "/signUp",
		Body:           `{"username": "test2username", "password": ""}`,
		Headers:        map[string]string{"Content-Type": "application/json"},
		ExpectedStatus: 400,
		ExpectedBody:   "Username and password are required",
	},
	{
		Name:           "SignUp with existing username",
		Method:         "POST",
		URL:            "/signUp",
		Body:           fmt.Sprintf(`{"username": "test%dusername", "password": "test1password", "email": "test@example.com"}`, randomTest),
		Headers:        map[string]string{"Content-Type": "application/json"},
		ExpectedStatus: 409,
		ExpectedBody:   "User already exists",
	},
}
