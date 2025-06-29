package handler

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"io"
	model "learning_go/internal/models"
	"log"
	"net/http"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

type compileBody struct {
	ID   string `json:"problemId"`
	Code string `json:"code"`
}

type compileRequest struct {
	Program   string          `json:"program"`
	FunName   string          `json:"funName"`
	TestCases [][]interface{} `json:"testCases"`
}

// GetFullCompile handles code compilation requests with caching
func GetFullCompile(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate HTTP method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var body compileBody

		// Read and parse request body
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		log.Printf("Compile endpoint received: ProblemID=%s, Code=%s", body.ID, body.Code)

		// Create hash of request body for caching
		bodyBytes, _ := json.Marshal(body)
		hash := sha256.Sum256(bodyBytes)
		hashStr := hex.EncodeToString(hash[:])

		// Check cache first
		if cached, exists := compileCache.Get(hashStr); exists {
			log.Printf("Cache hit for compile request: %s", hashStr)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(cached.StatusCode)
			json.NewEncoder(w).Encode(cached.ResponseBody)
			return
		}

		log.Printf("Cache miss for compile request: %s", hashStr)

		problemService := model.NewProblemService(db)
		problem, err := problemService.GetProblemByID(ctx, body.ID)

		if err != nil {
			log.Printf("Problem not found: %v", err)
			http.Error(w, "Problem not found", http.StatusNotFound)
			return
		}

		// Transform test cases to the expected format (inputs only)
		var transformedTestCases [][]interface{}
		for _, testCase := range problem.TestCases {
			// Parse input string into individual parameters
			inputs := strings.Fields(testCase.Input)
			var inputParams []interface{}
			for _, input := range inputs {
				if val, err := strconv.Atoi(input); err == nil {
					inputParams = append(inputParams, val)
				} else {
					inputParams = append(inputParams, input)
				}
			}

			// Create test case with just inputs (no expected output)
			transformedTestCases = append(transformedTestCases, inputParams)
		}

		// Create compile request body
		compileReq := compileRequest{
			Program:   body.Code,
			FunName:   problem.FunctionName,
			TestCases: transformedTestCases,
		}

		// Create HTTP client with TLS config to skip certificate verification
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		// Create request to compile service with the body
		compileReqBytes, err := json.Marshal(compileReq)
		if err != nil {
			http.Error(w, "Failed to marshal compile request", http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest("POST", "https://172.16.30.3:3001/performTestCases", bytes.NewBuffer(compileReqBytes))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Failed to create request",
			})
			return
		}

		// Print request details for debugging
		log.Printf("Sending compile request to service: %s", req.URL.String())
		log.Printf("Compile request body: %s", compileReqBytes)

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Compile service unavailable",
			})
			return
		}

		defer resp.Body.Close()

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Failed to read response",
			})
			return
		}

		var response struct {
			Results []struct {
				Output int    `json:"output"`
				Error  string `json:"error"`
				Line   int    `json:"line"`
				Column int    `json:"column"`
			} `json:"results"`
		}

		if err := json.Unmarshal(respBody, &response); err != nil {
			http.Error(w, "Failed to parse response", http.StatusInternalServerError)
			return
		}

		var structuredResponse model.CompileResponse
		// Actual response
		log.Printf("Compile service response: %+v", response)

		// Process all test results
		var results []model.CompileResults
		var firstError string
		var firstErrorLine, firstErrorColumn int
		hasError := false

		for i, result := range response.Results {
			// Check for compilation error
			if result.Error != "" {
				if !hasError {
					hasError = true
					firstError = result.Error
					firstErrorLine = result.Line
					firstErrorColumn = result.Column
				}
				results = append(results, model.CompileResults{
					Status:         "Failed",
					Output:         []int{result.Output},
					ExpectedOutput: []int{},
				})
				continue
			}

			// Compare with expected output
			expectedOutput := 0
			if i < len(problem.TestCases) {
				if val, err := strconv.Atoi(problem.TestCases[i].Output); err == nil {
					expectedOutput = val
				}
			}

			status := "Failed"
			if result.Output == expectedOutput {
				status = "Success"
			} else {
				hasError = true
			}

			results = append(results, model.CompileResults{
				Status:         status,
				Output:         []int{result.Output},
				ExpectedOutput: []int{expectedOutput},
			})
		}

		// Build response
		structuredResponse = model.CompileResponse{
			Result: results,
			Error:  firstError,
			Status: "Success",
			Line:   firstErrorLine,
			Column: firstErrorColumn,
		}

		if hasError {
			structuredResponse.Status = "Error"
		}

		// Cache successful responses
		if resp.StatusCode == http.StatusOK {
			compileCache.Set(hashStr, structuredResponse, resp.StatusCode)
			log.Printf("Cached compile response for request: %s", hashStr)
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)

		// Write response back to client
		json.NewEncoder(w).Encode(structuredResponse)
	}
}
