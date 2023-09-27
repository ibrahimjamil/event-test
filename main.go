package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

type CustomEvent struct {
	CustomEventName string `json:"customEventName"`
}

func handler(event events.CloudWatchEvent) (string, error) {
	var customEvent CustomEvent

	// Unmarshal the custom event data from the event detail
	err := json.Unmarshal([]byte(event.Detail), &customEvent)
	if err != nil {
		return "", err
	}

	// Access custom event properties
	eventName := customEvent.CustomEventName

	// Perform logic based on the custom event data
	result := fmt.Sprintf("Received custom event: Name=%s, Data=%s", eventName)
	fmt.Printf(result)

	return result, nil
}
