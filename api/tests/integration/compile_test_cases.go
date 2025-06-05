package integration

var Compile = []TestCase{
	{
		Name:           "Compile code with valid token",
		Method:         "POST",
		URL:            "/compile",
		Headers:        map[string]string{"Content-Type": "application/json", "Authorization": tokenString},
		Body:           `{code: "Int main() { return 0; }", language: "c", problemId: "6840ec83e844d5fee940c052"}`,
		ExpectedStatus: 200,
		ExpectedBody:   `[{"id":`,
	},
	{
		Name:           "Compile code with invalid token",
		Method:         "POST",
		URL:            "/compile",
		Headers:        map[string]string{"Content-Type": "application/json", "Authorization": badToken},
		Body:           `{code: "Int main() { return 0; }", language: "c", problemId: "6840ec83e844d5fee940c052"}`,
		ExpectedStatus: 401,
		ExpectedBody:   "Invalid token",
	},
}
