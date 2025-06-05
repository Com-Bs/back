package integration

var GetProblems = []TestCase{
	{
		Name:           "Get all problems with valid token",
		Method:         "GET",
		URL:            "/problems",
		Headers:        map[string]string{"Content-Type": "application/json", "Authorization": tokenString},
		ExpectedStatus: 200,
		ExpectedBody:   `[{"id":`,
	},
	{
		Name:           "Get all problems with invalid token",
		Method:         "GET",
		URL:            "/problems",
		Headers:        map[string]string{"Content-Type": "application/json", "Authorization": badToken},
		ExpectedStatus: 401,
		ExpectedBody:   "Invalid token",
	},
}

var GetProblemByID = []TestCase{
	{
		Name:           "Get specific problem with valid token",
		Method:         "GET",
		URL:            "/problems/6840ec83e844d5fee940c052",
		Headers:        map[string]string{"Content-Type": "application/json", "Authorization": tokenString},
		ExpectedStatus: 200,
		ExpectedBody:   `{"id":`,
	},
	{
		Name:           "Get specific problem with invalid token",
		Method:         "GET",
		URL:            "/problems/6840ec83e844d5fee940c052",
		Headers:        map[string]string{"Content-Type": "application/json", "Authorization": badToken},
		ExpectedStatus: 401,
		ExpectedBody:   "Invalid token",
	},
}
