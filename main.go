package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-stomp/stomp"
)

func main() {
	lambda.Start(handler)
}

type CloudWatchEvent struct {
	Version    string          `json:"version"`
	ID         string          `json:"id"`
	DetailType string          `json:"detail-type"`
	Source     string          `json:"source"`
	AccountID  string          `json:"account"`
	Time       time.Time       `json:"time"`
	Region     string          `json:"region"`
	Resources  []string        `json:"resources"`
	Detail     json.RawMessage `json:"detail"`
}

func extractJobNameFromEvent(resourceEvent string) (string, error) {
	jobName := strings.Split(resourceEvent, "/")
	return jobName[len(jobName)-1], nil
}

func handler(event events.CloudWatchEvent) (string, error) {
	var jobName string
	brokerEndpointIP := os.Getenv("MQ_ENDPOINT_IP")
	brokerUsername := os.Getenv("BROKER_USERNAME")
	brokerPassword := os.Getenv("BROKER_PASSWORD")
	brokerEndpointIP = strings.TrimPrefix(brokerEndpointIP, "stomp+ssl://")

	fmt.Println(event)
	if len(event.Resources) > 0 {
		resourceEvent := event.Resources[0]
		value, err := extractJobNameFromEvent(resourceEvent)
		if err != nil {
			return fmt.Sprintf("Failed to extract job name from event: %v", err), err
		}
		jobName = value
		fmt.Printf("Received job name=%s", jobName)
	} else {
		fmt.Println("No rule resources found in the CloudWatch event.")
	}

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
		return fmt.Sprintf("Failed to connect to the broker: %v", err), err
	}
	defer conn.Disconnect()

	fmt.Print("connection established to broker....")

	// Send a message to a queue on the broker
	queueName := "Demo-Queue"
	message := jobName
	err = conn.Send(
		queueName,
		"text/plain",
		[]byte(message),
		nil,
	)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return fmt.Sprintf("Failed to send message: %v", err), err
	}

	log.Printf("Message sent to the queue: %s", message)

	// Subscribe to a queue on the broker
	sub, err := conn.Subscribe(queueName, stomp.AckAuto)
	if err != nil {
		log.Printf("Failed to subscribe to the queue: %v", err)
		return fmt.Sprintf("Failed to subscribe to the queue: %v", err), err
	}
	defer sub.Unsubscribe()

	fmt.Print("Connection established for recieving, waiting for messages...\n")

	// Listen for and process incoming messages
	var messageBody string
	for {
		msg := <-sub.C
		if msg.Err != nil {
			log.Printf("Failed to receive message: %v", msg.Err)
			return fmt.Sprintf("Failed to receive message: %v", msg.Err), msg.Err
		}

		// Process the received message (you can modify this part as needed)
		messageBody = string(msg.Body)
		log.Printf("Received message from the queue: %s", messageBody)
		break
	}

	return fmt.Sprintf("Received custom event: Name=%s", jobName), nil
}
