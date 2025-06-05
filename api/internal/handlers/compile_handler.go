package handler

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	Program   string      `json:"program"`
	FunName   string      `json:"funName"`
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
			w.Write(cached.ResponseBody)
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

		req, err := http.NewRequest("POST", "https://10.49.12.48:3001/runCompile", bytes.NewBuffer(compileReqBytes))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		// Print request details for debugging
		log.Printf("Sending compile request to service: %s", req.URL.String())
		log.Printf("Compile request body: %s", compileReqBytes)

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Compile service unavailable", http.StatusServiceUnavailable)
			return
		}

		defer resp.Body.Close()

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response", http.StatusInternalServerError)
			return
		}

		var response struct {
			Outputs []int  `json:"outputs"`
			Error   string `json:"error"`
			Message string `json:"message"`
			Line    int    `json:"line"`
			Column  int    `json:"column"`
		}

		if err := json.Unmarshal(respBody, &response); err != nil {
			http.Error(w, "Failed to parse response", http.StatusInternalServerError)
			return
		}

		var structuredResponse model.CompileResponse
		if response.Error != "" {
			log.Printf("Compile service returned an error: %s", response.Error)
			errorMessage := response.Error
			if response.Message != "" {
				errorMessage = fmt.Sprintf("%s: %s", response.Error, response.Message)
			}
			structuredResponse = model.CompileResponse{
				Error:  errorMessage,
				Line:   response.Line,
				Column: response.Column,
				Status: "Error",
			}
		} else {
			// For now, create a simple success response
			// TODO: Update GenerateResponse to handle TestCase structures
			structuredResponse = model.CompileResponse{
				Error:  "",
				Status: "Success",
				Line:   response.Line,
				Column: response.Column,
			}
		}

		// Cache successful responses
		if resp.StatusCode == http.StatusOK {
			compileCache.Set(hashStr, respBody, resp.StatusCode)
			log.Printf("Cached compile response for request: %s", hashStr)
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)

		// Write response back to client
		log.Printf("Compile endpoint returning: %+v", structuredResponse)
		json.NewEncoder(w).Encode(structuredResponse)
	}
}
