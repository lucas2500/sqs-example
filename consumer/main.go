package main

import (
	"log"
	"os"
	"sqs-example/queue"
	"sqs-example/util"
	"sync"

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

	wg := sync.WaitGroup{}

	messages := make(chan []types.Message)

	QueueName := os.Getenv("SQS_QUEUE_NAME")

	wg.Add(1)

	go func() {

		defer wg.Done()
		queue.DequeueMessage(QueueName, messages)
	}()

	for message := range messages {

		log.Println("- Received message:", *message[0].Body)

		queue.RemoveMessage(QueueName, *message[0].ReceiptHandle)
	}

	wg.Wait()
}
