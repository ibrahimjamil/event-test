package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(event interface{}) (string, error) {
	// Type assertion to convert event into a struct type
	customEvent, ok := event.(struct {
		CustomEventName string `json:"customEventName"`
	})

	if !ok {
		return "", fmt.Errorf("Failed to assert event to the expected type")
	}

	// Access the customEventName field
	eventName := customEvent.CustomEventName

	// Use the custom event data in your logic
	result := fmt.Sprintf("Received custom event: Name=%s", eventName)

	return result, nil
}
