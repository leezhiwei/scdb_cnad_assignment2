package main

// Import packages
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Request structure for Med42 API
type Med42Request struct {
	Question string `json:"question"`
}

// Response structure for Med42 API
type Med42Response struct {
	Answer string `json:"answer"`
}

func main() {
	// Static testing question for med llm
	question := "If I fall, what should I do?"

	// Med42 API endpoint and API key
	apiURL := "https://api.med42.com/v1/ask"
	// Obtain API key from environment
	apiKey := os.Getenv("MED42_API_KEY")

	// Prepare the request payload
	requestBody, err := json.Marshal(Med42Request{Question: question})
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// Create an HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("Error creating HTTP request: %v\n", err)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d\n", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Response: %s\n", body)
		return
	}

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	var med42Response Med42Response
	err = json.Unmarshal(body, &med42Response)
	if err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return
	}

	// Testing response
	fmt.Printf("Q: %s\nA: %s\n", question, med42Response.Answer)
}
