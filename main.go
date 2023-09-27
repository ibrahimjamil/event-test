package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(event map[string]interface{}) (string, error) {
	// Access the "EventName" field from the map
	eventName, ok := event["jobName"].(string)
	if !ok {
		return "", fmt.Errorf("Failed to extract EventName from event")
	}

	// Use the eventName in your logic
	result := fmt.Sprintf("Received custom event: Name=%s", eventName)

	fmt.Printf("Received customize event: Name=%s", eventName)

	return result, nil
}
