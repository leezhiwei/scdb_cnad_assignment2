package main

import (
	"fmt"
	"log"
	"context"
	"github.com/smsapi/smsapi-go/smsapi"

func sendSMS(to string, code string) error {
	// credentials
	accessToken := ""

	client = smsapi.NewInternationalClient(accessToken, nil)

	// Create the message
	body := fmt.Sprintf("Your login code is: %s", code)

	// Send the SMS
	result, err := client.Sms.Send(context.Background(), to, body, "")
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	log.Printf("SMS sent successfully")
	fmt.Println("Sent messages count", result.Count)

	for _, sms := range result.Collection {
		fmt.Println(sms.Id, sms.Status, sms.Points)
	}
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
