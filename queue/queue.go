package queue

import (
	"context"
	"log"
	"sqs-example/util"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSClient interface {
	GetQueueUrl(ctx context.Context, params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)

	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)

	SendMessage(ctx context.Context, params *sqs.SendMessageInput,
		optsFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

func GetQueueURL(queue string) *sqs.GetQueueUrlOutput {

	cfg := util.Cfg

	client := sqs.NewFromConfig(cfg)

	QueueInput := &sqs.GetQueueUrlInput{
		QueueName: &queue,
	}

	url, err := client.GetQueueUrl(context.TODO(), QueueInput)

	if err != nil {
		log.Fatal("- Unable to get queue URL")
	}

	return url
}

func GetMessages(ctx context.Context, s SQSClient, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return s.ReceiveMessage(ctx, input)
}

func RemoveMessage(queue string, ReceiptHandle string) {

	cfg := util.Cfg

	client := sqs.NewFromConfig(cfg)

	url := GetQueueURL(queue)

	DeleteMessageInput := &sqs.DeleteMessageInput{
		QueueUrl:      url.QueueUrl,
		ReceiptHandle: &ReceiptHandle,
	}

	_, err := client.DeleteMessage(context.Background(), DeleteMessageInput)

	if err != nil {
		log.Fatal("- Unable to delete message from the queue", queue)
	}

	log.Println("- Message deleted successfully")
}

func QueueMessage(MessageObject *sqs.SendMessageInput) {

	cfg := util.Cfg

	client := sqs.NewFromConfig(cfg)

	_, err := client.SendMessage(context.Background(), MessageObject)

	if err != nil {
		log.Fatal("- Unable to queue message")
	}

	log.Println("- Message queued successfully")
}

func DequeueMessage(queue string, messages chan []types.Message) {

	forever := make(chan bool)

	cfg := util.Cfg

	client := sqs.NewFromConfig(cfg)

	url := GetQueueURL(queue)

	go func() {

		for {

			MessageInput := &sqs.ReceiveMessageInput{
				MessageAttributeNames: []string{
					string(types.QueueAttributeNameAll),
				},
				QueueUrl:            url.QueueUrl,
				MaxNumberOfMessages: 1,
				VisibilityTimeout:   5,
			}

			MessageResult, err := GetMessages(context.TODO(), client, MessageInput)

			if err != nil {
				log.Fatal("- Unable to get messages")
			}

			if MessageResult.Messages != nil {
				messages <- MessageResult.Messages
			}
		}
	}()

	log.Println("- Worker up and running...")

	<-forever
}
