package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type MyMessage struct {
	Message string `json:"message"`
}

func HandleRequest(ctx context.Context, msg MyMessage) (string, error) {
	revMessage := Reverse(msg.Message)
	SendSQSMessage(revMessage)
	return fmt.Sprintf("%s!", revMessage), nil
}

func SQSHandler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
	}

	return nil
}

func Reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

func SendSQSMessage(s string) error {
	sess := session.New(&aws.Config{
		Region: aws.String("eu-central-1"),
	})

	svc := sqs.New(sess)
	qURL := "https://sqs.eu-central-1.amazonaws.com/445174708070/demo_queue"

	result, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Title": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("The Whistler"),
			},
			"Author": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("John Grisham"),
			},
			"WeeksOn": &sqs.MessageAttributeValue{
				DataType:    aws.String("Number"),
				StringValue: aws.String("6"),
			},
		},
		MessageBody: aws.String("Information about current NY Times fiction bestseller for week of 12/11/2016."),
		QueueUrl:    &qURL,
	})

	if err != nil {
		fmt.Println("Error", err)
		return err
	}

	fmt.Println("Success", *result.MessageId)

	return nil
}

func recvSQSMessage() (string, error) {
	sess := session.New(&aws.Config{
		Region: aws.String("eu-central-1"),
	})

	service := sqs.New(sess)
	qURL := "https://sqs.eu-central-1.amazonaws.com/445174708070/Demo"

	receive_params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(qURL),
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(30),
		WaitTimeSeconds:     aws.Int64(1),
	}
	receive_resp, err := service.ReceiveMessage(receive_params)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("[Receive message] \n%v \n\n", receive_resp)

	return "ok", nil
}

func main() {
	lambda.Start(SQSHandler)
}
