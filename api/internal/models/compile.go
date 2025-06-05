package model

type CompileResponse struct {
	Result []CompileResults `json:"result,omitempty"` // Optional field for compile results
	Error  string           `json:"error"`
	Status string           `json:"status,omitempty"` // Optional field for status
	Line   int              `json:"line,omitempty"`   // Optional field for line number
	Column int              `json:"column,omitempty"` // Optional field for column number
}

type CompileResults struct {
	Status         string `json:"status"`
	Output         []int  `json:"output"`
	ExpectedOutput []int  `json:"expectedOutput"`
}

func GenerateResponse(output []int, expectedOutput []int) CompileResponse {
	var results []CompileResults
	statusStr := "Success"

	for i, out := range output {
		var respObj CompileResults
		if out == expectedOutput[i] {
			respObj = CompileResults{
				Status:         "Success",
				Output:         []int{out},
				ExpectedOutput: []int{expectedOutput[i]},
			}
		} else {
			respObj = CompileResults{
				Status:         "Failed",
				Output:         []int{out},
				ExpectedOutput: []int{expectedOutput[i]},
			}

			statusStr = "Error"
		}
		results = append(results, respObj)
	}

	return CompileResponse{
		Result: results,
		Error:  "",
		Status: statusStr,
		Line:   0, // Default value, can be set if needed
		Column: 0, // Default value, can be set if needed
	}
}
