package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func handleCloudWatchEvent(event *events.CloudWatchEvent) (interface{}, error) {
	// Handle CloudWatch Event (event logs) here
	// Access event.Detail and other properties as needed
	fmt.Printf(string(event.Detail))
	return string(event.Detail), nil
}
