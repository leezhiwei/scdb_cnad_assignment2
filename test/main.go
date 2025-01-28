package main

import (
	"fmt"
	"log"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func sendSMS(to string, code string) error {
	// Twilio credentials
	accountSid := ""
	authToken := ""
	from := "" // Your Twilio phone number

	// Create a Twilio client
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	// Create the message
	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(fmt.Sprintf("Your login code is: %s", code))

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
	toPhoneNumber := "+65 85178498" // User's phone number
	loginCode := "123456"           // Randomly generated login code

	if err := sendSMS(toPhoneNumber, loginCode); err != nil {
		log.Fatalf("Error sending SMS: %v", err)
	}
}
