package main

// Import packages
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// curl.exe -X POST http://localhost:5000/api/v1/medqna -H "Content-Type: application/json" -d "{\"question\": \"If i fall, what should i do?\"}"

// Request and Response struct
type QnARequest struct {
	Question string `json:"question"`
}
type QnAResponse struct {
	Answer string `json:"answer"`
}

// User request struct
type UserMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AI request and response struct
type AIRequest struct {
	Model    string        `json:"model"`
	Messages []UserMessage `json:"messages"`
}
type AIResponse struct {
	Choices []struct {
		Message UserMessage `json:"message"`
	} `json:"choices"`
}

// Handler function for path /api/v1/medqna
func medicalqna(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON body from the client
	var req QnARequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	// Display received question on go terminal
	fmt.Printf("Received Question: %s\n", req.Question)

	// Construct the AI model request payload
	aiReq := AIRequest{
		Model: "m42-health/Llama3-Med42-8B",
		Messages: []UserMessage{
			{
				// Set behaviour of the model
				Role:    "system",
				Content: "Always answer as helpfully as possible, while being safe. Your answers should not include any harmful, unethical, racist, sexist, toxic, dangerous, or illegal content. Please ensure that your responses are socially unbiased and positive in nature. If a question does not make any sense, or is not factually coherent, explain why instead of answering something not correct. If you don’t know the answer to a question, please don’t share false information.",
			},
			{
				// User question
				Role:    "user",
				Content: req.Question,
			},
		},
	}

	// Convert AI request to JSON
	aiReqBody, err := json.Marshal(aiReq)
	if err != nil {
		http.Error(w, "Error encoding AI request", http.StatusInternalServerError)
		return
	}

	// POST request to the Med42 Llama3 model endpoint
	aiEndpoint := "http://192.168.2.108:8000/v1/chat/completions"
	resp, err := http.Post(aiEndpoint, "application/json", bytes.NewBuffer(aiReqBody))
	if err != nil {
		http.Error(w, "Error connecting to AI model", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read Med42 Llama3 model response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading AI response", http.StatusInternalServerError)
		return
	}

	// Parse response to JSON
	var aiResp AIResponse
	err = json.Unmarshal(body, &aiResp)
	if err != nil {
		http.Error(w, "Error parsing AI response", http.StatusInternalServerError)
		return
	}

	// Set default response to no response
	aiAnswer := "No response received"
	// Check if there is a response from Med42 Llama3 model
	if len(aiResp.Choices) > 0 {
		aiAnswer = aiResp.Choices[0].Message.Content
	}

	// Return AI-generated response to client
	response := QnAResponse{Answer: aiAnswer}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	// Display Med42 Llama3 response on go terminal
	fmt.Printf("Med42 Llama3 Response: %s\n", response.Answer)
}

// // Display searched course information function via browser or cmd
// func medicalqna(w http.ResponseWriter, r *http.Request) {
// 	// Check if the request method is POST
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Variable to store the created or updated course information
// 	var tempQna Med42Question
// 	// Parse the JSON body to get course details and add into tempCourse variable
// 	errorStatus := json.NewDecoder(r.Body).Decode(&tempQna)
// 	if errorStatus != nil {
// 		// Write to client with status internal server error (500) from JSON serialization
// 		w.WriteHeader(http.StatusInternalServerError)
// 		fmt.Fprintf(w, "ERROR from JSON serialization")
// 		return
// 	}
// 	// Print the decoded JSON object
// 	fmt.Printf("Decoded JSON: %+v\n", tempQna)

// 	// Secure HTTP headers
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("X-Content-Type-Options", "nosniff")
// 	w.Header().Set("X-Frame-Options", "DENY")
// 	w.Header().Set("X-XSS-Protection", "1; mode=block")

// 	// Encode the response and handle errors
// 	if err := json.NewEncoder(w).Encode(tempQna); err != nil {
// 		log.Println("JSON encoding error:", err)
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 	}
// }

func main() {
	// Mux router for routing HTTP request
	router := mux.NewRouter()
	// Mux router to map path to different functions
	router.HandleFunc("/api/v1/medqna", func(w http.ResponseWriter, r *http.Request) {
		medicalqna(w, r)
	})

	// Listen at port 5000
	//fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
