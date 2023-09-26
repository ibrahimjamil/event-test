package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-stomp/stomp/v3"
)

func main() {
	lambda.Start(Handler{})
}

func handleAPIGatewayEvent(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get the broker endpoint
	brokerEndpointIP := os.Getenv("MQ_ENDPOINT_IP")
	brokerUsername := os.Getenv("BROKER_USERNAME")
	brokerPassword := os.Getenv("BROKER_PASSWORD")
	brokerEndpointIP = strings.TrimPrefix(brokerEndpointIP, "stomp+ssl://")

	// Create a tls dial and stomp connect to broker
	netConn, err := tls.Dial("tcp", brokerEndpointIP, &tls.Config{})
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer netConn.Close()

	conn, err := stomp.Connect(netConn,
		stomp.ConnOpt.Login(brokerUsername, brokerPassword))
	if err != nil {
		log.Printf("Failed to connect to the broker: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}
	defer conn.Disconnect()

	fmt.Print("connection established")

	// Send a message to a queue on the broker
	queueName := "Demo-Queue"
	message := request.Body
	err = conn.Send(
		queueName,
		"text/plain",
		[]byte(message),
		nil,
	)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	log.Printf("Message sent to the queue: %s", message)

	// Subscribe to a queue on the broker
	sub, err := conn.Subscribe(queueName, stomp.AckAuto)
	if err != nil {
		log.Printf("Failed to subscribe to the queue: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}
	defer sub.Unsubscribe()

	fmt.Print("Connection established, waiting for messages...\n")

	// Listen for and process incoming messages
	var messageBody string
	for {
		msg := <-sub.C
		if msg.Err != nil {
			log.Printf("Failed to receive message: %v", msg.Err)
			return events.APIGatewayProxyResponse{StatusCode: 500}, msg.Err
		}

		// Process the received message (you can modify this part as needed)
		messageBody = string(msg.Body)
		log.Printf("Received message from the queue: %s", messageBody)
		break
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("Message sent: %s and recieved also %s", "done", messageBody),
	}
	return response, nil
}

func handleCloudWatchEvent(event events.CloudWatchEvent) (interface{}, error) {
	// Handle CloudWatch Event (event logs) here
	// Access event.Detail and other properties as needed
	fmt.Printf(string(event.Detail))
	return string(event.Detail), nil
}

type Record struct {
	EventVersion                   string `json:"EventVersion"`
	EventSubscriptionArn           string `json:"EventSubscriptionArn"`
	EventSource                    string `json:"EventSource"`
	events.APIGatewayProxyResponse `json:",omitempty"`
}

type Event struct {
	Records Record `json:"Records"`
}

type Response struct {
	events.APIGatewayProxyResponse `json:",omitempty"`
	//other response type
}

type CustomEvent struct {
	//Give EventSource "aws:manual" value to determine event is manual
	EventSource string `json:"eventSource"`

	//Other CustomEvent Attributes
}

type Handler struct {
	//add global variables or context information that your handler may need
}

func (h Handler) Invoke(ctx context.Context, data []byte) ([]byte, error) {
	// fmt.Printf("event-source", event.Records.EventSource)
	// fmt.Printf(event.Records[0].EventSource)
	// if apiGatewayEvent, ok := request.(events.APIGatewayProxyRequest); ok {
	// 	return handleAPIGatewayEvent(apiGatewayEvent)
	// }

	// if cloudWatchEvent, ok := request.(events.CloudWatchEvent); ok {
	// 	return handleCloudWatchEvent(cloudWatchEvent)
	// }

	// var response Response
	// switch {
	// case reflect.DeepEqual(event.APIGatewayProxyRequest, events.APIGatewayProxyRequest{}):
	// 	response.APIGatewayProxyResponse, _ = handleAPIGatewayEvent(event.APIGatewayProxyRequest)
	// 	//another case for a event handler
	// 	fmt.Printf("hello")
	// }

	//for demonstration purposes, not the best way to handle
	fmt.Printf("hello")
	apiGatewayEvent := new(events.APIGatewayProxyRequest)
	if err := json.Unmarshal(data, apiGatewayEvent); err != nil {
		log.Println("Not a api gateway event")
	}
	snsEvent := new(events.SNSEvent)
	if err := json.Unmarshal(data, snsEvent); err != nil {
		log.Println("Not a sns event")
	}
	return nil, nil
}
