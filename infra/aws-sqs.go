package infra

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Queue struct {
	client *sqs.SQS
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Setup(region, clientId, clientSecret string) error {
	c := credentials.NewStaticCredentials(clientId, clientSecret, "")
	sess, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(region),
			Credentials: c,
		},
	)
	if err != nil {
		fmt.Println("Error creating session:", err)
		return err
	}
	client := sqs.New(sess)
	q.client = client
	return nil
}

func (q *Queue) SendMessage(queueURL string, message string) error {
	fmt.Println("========================")
	fmt.Println("SendMessage:", queueURL)
	fmt.Println("========================")
	req, resp := q.client.SendMessageRequest(&sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String(queueURL),
	})
	err := req.Send()
	if err != nil {
		fmt.Println("Error sending message:", err)
		return err
	}
	fmt.Println("Resp: ", resp)
	return nil
}

func (q *Queue) ConsumeMessages(queueURL, message string) error {
	fmt.Println("========================")
	fmt.Println("SendMessage:", queueURL)
	fmt.Println("========================")
	req, resp := q.client.SendMessageRequest(&sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String(queueURL),
	})
	err := req.Send()
	if err != nil {
		fmt.Println("Error sending message:", err)
		return err
	}
	fmt.Println("Resp: ", resp)
	return nil
}
