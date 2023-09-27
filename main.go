package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(event interface{}) (string, error) {
	fmt.Println(event)

	return "", nil
}
