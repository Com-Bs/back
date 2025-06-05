package handler

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	model "learning_go/internal/models"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type compileBody struct {
	ID   string `json:"problemId"`
	Code string `json:"code"`
}

type compileRequest struct {
	Code      string `json:"code"`
	TestCases []int  `json:"cases"`
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

		problemService := model.NewProblemService(db)
		problem, err := problemService.GetProblemByID(ctx, body.ID)

		if err != nil {
			log.Printf("Problem not found: %v", err)
			http.Error(w, "Problem not found", http.StatusNotFound)
			return
		}

		// Create compile request body
		compileReq := compileRequest{
			Code:      body.Code,
			TestCases: problem.TestCases,
		}

		// Create HTTP client
		client := &http.Client{}

		log.Printf("Cache miss for compile request: %s", hashStr)

		// Create request to compile service with the body
		compileReqBytes, err := json.Marshal(compileReq)
		if err != nil {
			http.Error(w, "Failed to marshal compile request", http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest("POST", "http://10.49.12.48:3001/runCompile", bytes.NewBuffer(compileReqBytes))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Failed to send request", http.StatusInternalServerError)
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
			Output []int `json:"output"`
			Error  int   `json:"error"`
			Line   int   `json:"line"`
			Column int   `json:"column"`
		}

		if err := json.Unmarshal(respBody, &response); err != nil {
			http.Error(w, "Failed to parse response", http.StatusInternalServerError)
			return
		}

		var structuredResponse model.CompileResponse
		if response.Error != 0 {
			log.Printf("Compile service returned an error: %d", response.Error)
			structuredResponse = model.CompileResponse{
				Error:  fmt.Sprintf("Compilation error: %d", response.Error),
				Line:   response.Line,
				Column: response.Column,
				Status: "Error",
			}
		} else {
			structuredResponse = model.GenerateResponse(response.Output, problem.TestCases)
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
		json.NewEncoder(w).Encode(structuredResponse)
	}
}
