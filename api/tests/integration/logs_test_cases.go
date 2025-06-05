package integration

import (
	"fmt"
	"learning_go/internal/auth"
)

func GetToken() (string, string) {
	jwt, err := auth.CreateToken("testuser")

	if err != nil {
		panic(fmt.Sprintf("Error creating token: %v", err))
	}

	tokenString := fmt.Sprintf("Bearer %s", jwt)
	modifiedToken := tokenString[:len(tokenString)-1] + "y" // Modify the last character to 'y' for testing purposes

	return tokenString, modifiedToken
}

var tokenString, badToken = GetToken()

var GetLogs = []TestCase{
	{
		Name:           "Get logs with valid token",
		Method:         "GET",
		URL:            "/logs",
		Headers:        map[string]string{"Content-Type": "application/json", "Authorization": tokenString},
		ExpectedStatus: 200,
		ExpectedBody:   `[{"ID":`,
	},
	{
		Name:           "Get logs with invalid token",
		Method:         "GET",
		URL:            "/logs",
		Headers:        map[string]string{"Content-Type": "application/json", "Authorization": badToken},
		ExpectedStatus: 401,
		ExpectedBody:   "Invalid token",
	},
}
