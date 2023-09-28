package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/go-stomp/stomp"
)

func main() {
	lambda.Start(handler)
}

func getSecret(client *ssm.SSM, paramName string) (string, error) {
	result, err := client.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(paramName),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}
	return *result.Parameter.Value, nil
}

func handler(event map[string]interface{}) (string, error) {
	// Access the "EventName" field from the map
	eventName, ok := event["jobName"].(string)
	if !ok {
		return "", fmt.Errorf("Failed to extract EventName from event")
	}

	// Use the eventName in your logic
	// result := fmt.Sprintf("Received custom event: Name=%s", eventName)

	// fmt.Printf("Received customize event: Name=%s", eventName)

	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1"), // Replace with your AWS region
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	// Create an AWS Systems Manager service client
	ssmClient := ssm.New(sess)
	fmt.Printf("Failed to get broker password: %v", err)

	// Retrieve secrets from Parameter Store
	brokerUsername, err := getSecret(ssmClient, "/event-scheduling/broker_username")
	if err != nil {
		fmt.Printf("Failed to get broker username: %v", err)
		return fmt.Sprintf("Failed to get broker username: %v", err), err
	}
	brokerPassword, err := getSecret(ssmClient, "/event-scheduling/broker_password")
	if err != nil {
		fmt.Printf("Failed to get broker password: %v", err)
		return fmt.Sprintf("Failed to get broker password: %v", err), err
	}

	// Get the broker endpoint
	brokerEndpointIP := os.Getenv("MQ_ENDPOINT_IP")
	brokerEndpointIP = strings.TrimPrefix(brokerEndpointIP, "stomp+ssl://")

	// Create a tls dial and stomp connect to broker
	netConn, err := tls.Dial("tcp", brokerEndpointIP, &tls.Config{})
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer netConn.Close()
	fmt.Printf("Received customize event: Name=%s", eventName)
	conn, err := stomp.Connect(netConn,
		stomp.ConnOpt.Login(brokerUsername, brokerPassword))
	if err != nil {
		log.Printf("Failed to connect to the broker: %v", err)
		return fmt.Sprintf("Failed to connect to the broker: %v", err), err
	}
	defer conn.Disconnect()

	fmt.Print("connection established")

	// Send a message to a queue on the broker
	queueName := "Demo-Queue"
	message := eventName
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

	fmt.Print("Connection established, waiting for messages...\n")

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

	return fmt.Sprintf("Received custom event: Name=%s", eventName), nil
}
