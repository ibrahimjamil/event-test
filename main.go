package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(event events.CloudWatchEvent) (string, error) {
	// Handle CloudWatch Event (event logs) here
	// Access event.Detail and other properties as needed
	return fmt.Sprint("event", string(event.Source)), nil
}
