package main

import (
	"log"
	"os"
	"sqs-example/queue"
	"sqs-example/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/joho/godotenv"
)

func init() {

	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("- Unable to load .env")
	}

	util.LoadAWSConfig()
}

func main() {

	QueueName := os.Getenv("SQS_QUEUE_NAME")

	url := queue.GetQueueURL(QueueName)

	MessageInput := &sqs.SendMessageInput{
		DelaySeconds: 5,
		MessageAttributes: map[string]types.MessageAttributeValue{

			"Name": {
				DataType:    aws.String("String"),
				StringValue: aws.String("Lucas Barbosa"),
			},

			"Age": {
				DataType:    aws.String("Number"),
				StringValue: aws.String("24"),
			},
			"Position": {
				DataType:    aws.String("String"),
				StringValue: aws.String("Backend engineer"),
			},
		},
		MessageBody: aws.String("Some informations about me"),
		QueueUrl:    url.QueueUrl,
	}

	queue.QueueMessage(MessageInput)
}
