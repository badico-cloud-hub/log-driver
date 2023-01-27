package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	LogQ   amqp.Queue
	EventQ amqp.Queue
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Setup(username, password, serverUrl string) error {
	uri := fmt.Sprintf("amqps://%s:%s@%s", username, password, serverUrl)
	fmt.Println(uri)
	conn, err := amqp.Dial(uri)
	if err != nil {
		log.Panicf("%s: %s", "Failed to connect to RabbitMQ", err)
	}
	// defer conn.Close()

	q.conn = conn
	ch, err := q.conn.Channel()
	q.ch = ch
	logQ, err := ch.QueueDeclare(
		"LogMessages", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	q.LogQ = logQ

	// err check

	eventQ, err := ch.QueueDeclare(
		"EventMessages", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)

	// err check

	q.EventQ = eventQ
	return nil
}

func (q *Queue) SendMessage(nameq string, data interface{}) error {
	fmt.Println("========================")
	fmt.Println("SendMessage to ", nameq)
	fmt.Println("========================")
	message, err := json.Marshal(data)
	if err != nil {
		fmt.Println("============================================")
		fmt.Println("WARN: CLOUDLOG ENGINE NOT WORKING")
		fmt.Println("============================================")
		fmt.Println(err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	defer cancel()
	err = q.ch.PublishWithContext(ctx,
		"",    // exchange
		nameq, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		})

	if err != nil {
		fmt.Println("Error sending message:", err)
		return err
	}
	return nil
}

// func (q *Queue) ConsumeMessages(queueURL string) error {
// 	fmt.Println("========================")
// 	fmt.Println("SendMessage:", queueURL)
// 	fmt.Println("========================")
// 	req, resp := q.client.SendMessageRequest(&sqs.SendMessageInput{
// 		MessageBody: aws.String(message),
// 		QueueUrl:    aws.String(queueURL),
// 	})
// 	err := req.Send()
// 	if err != nil {
// 		fmt.Println("Error sending message:", err)
// 		return err
// 	}
// 	fmt.Println("Resp: ", resp)
// 	return nil
// }
