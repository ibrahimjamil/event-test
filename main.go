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

func handler(request interface{}) (interface{}, error) {
	// Try to unmarshal the request into an APIGatewayProxyRequest
	if gatewayRequestData, err := json.Marshal(request); err == nil {
		gatewayRequest := &events.APIGatewayProxyRequest{}
		if err := json.Unmarshal(gatewayRequestData, gatewayRequest); err == nil {
			if gatewayRequest.HTTPMethod == "POST" ||
				gatewayRequest.HTTPMethod == "PUT" ||
				gatewayRequest.HTTPMethod == "DELETE" ||
				gatewayRequest.HTTPMethod == "GET" {
				return handleAPIGatewayEvent(gatewayRequest)
			}
		}
	}

	// Try to unmarshal the request into a CloudWatchEvent
	if cloudWatchEventData, err := json.Marshal(request); err == nil {
		cloudWatchEvent := &events.CloudWatchEvent{}
		if err := json.Unmarshal(cloudWatchEventData, cloudWatchEvent); err == nil {
			if len(cloudWatchEvent.Detail) > 0 {
				fmt.Println("Received CloudWatch event with non-empty detail")
				return handleCloudWatchEvent(cloudWatchEvent)
			}
		}
	}

	return events.APIGatewayProxyResponse{StatusCode: 500}, nil
}
