package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

const debugMode = true // Set to true to use a fixed number

func generateRandomNumber() string {
	if debugMode {
		return "123456" // Fixed number for debugging
	}
	time.Now().UnixNano()
	return fmt.Sprintf("%06d", rand.Intn(999999-1)+1)
}

func sendSMS(to string, code string) error {
	// Twilio credentials
	accountSid := ""
	authToken := ""
	from := "+18454705971" // Your Twilio phone number

	// Create a Twilio client
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	// Create the message
	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(fmt.Sprintf("Your login code for SDCB is: %s", code))

	// Send the SMS
	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	log.Printf("SMS sent successfully: SID=%s", *resp.Sid)
	return nil
}

func main() {
	// Example usage
	toPhoneNumber := "+65 91542235"     // User's phone number
	loginCode := generateRandomNumber() // Randomly generated login code
	fmt.Println(loginCode)
	if !debugMode {
		if err := sendSMS(toPhoneNumber, loginCode); err != nil {
			log.Fatalf("Error sending SMS: %v", err)
		}
	}
}
