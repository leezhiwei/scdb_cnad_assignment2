package main

// Import packages
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	CORShandler "github.com/leezhiwei/common/CORSHandler"
	"github.com/leezhiwei/common/mainhandler"
	"github.com/leezhiwei/common/ping"
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

// Med42 request and response struct
type Med42Request struct {
	Model    string        `json:"model"`
	Messages []UserMessage `json:"messages"`
}
type Med42Response struct {
	Choices []struct {
		Message UserMessage `json:"message"`
	} `json:"choices"`
}

// Handler function for path /api/v1/medqna
func medicalqna(w http.ResponseWriter, r *http.Request) {
	// CORS settings
	var preflight bool = CORShandler.SetCORSHeaders(&w, r)
	if preflight {
		return
	}
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

	// Construct the Med42 Llama3 model request payload
	med42Req := Med42Request{
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

	// Convert Med42 Llama3 request to JSON
	med42ReqBody, err := json.Marshal(med42Req)
	if err != nil {
		http.Error(w, "Error encoding Med42 Llama3 request", http.StatusInternalServerError)
		return
	}

	// POST request to the Med42 Llama3 model endpoint
	med42Endpoint := config.API.Med42Endpoint
	resp, err := http.Post(med42Endpoint, "application/json", bytes.NewBuffer(med42ReqBody))
	if err != nil {
		http.Error(w, "Error connecting to Med42 Llama3 model", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read Med42 Llama3 model response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading Med42 Llama3 response", http.StatusInternalServerError)
		return
	}

	// Parse response to JSON
	var modelResp Med42Response
	err = json.Unmarshal(body, &modelResp)
	if err != nil {
		http.Error(w, "Error parsing Med42 Llama3 response", http.StatusInternalServerError)
		return
	}

	// Set default response to no response
	modelAnswer := "No response received"
	// Check if there is a response from Med42 Llama3 model
	if len(modelResp.Choices) > 0 {
		modelAnswer = modelResp.Choices[0].Message.Content
	}

	// Return Med42 Llama3 response to client
	response := QnAResponse{Answer: modelAnswer}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	// Display Med42 Llama3 response on go terminal
	fmt.Printf("Med42 Llama3 Response: %s\n", response.Answer)
}

type Config struct {
	API struct {
		Med42Endpoint string `json:"med42Endpoint"`
	} `json:"api"`
}

func GetConfig() Config {
	configFile, err := os.Open("./config.json")
	defer configFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

var config Config

func main() {
	var prefix string = "/api/v1/medqna"
	// Mux router for routing HTTP request
	router := mux.NewRouter()

	// Mux router pinghandler for service heartbeat
	router.HandleFunc(fmt.Sprintf("%s/ping", prefix), ping.PingHandler).Methods("GET", "OPTIONS")

	// Mux router to map path to different functions
	router.HandleFunc(fmt.Sprintf("%s/chat", prefix), func(w http.ResponseWriter, r *http.Request) {
		medicalqna(w, r)
	})
	// Logging IPs
	router.Use(mainhandler.LogReq)

	// Listen at port 5000
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServeTLS(":5000", "../certs/server.cert", "../certs/server.key", router))
}
