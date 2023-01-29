package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	LogQ   amqp.Queue
	EventQ amqp.Queue
}

type callbackFunc func(msg amqp.Delivery)

func NewQueue() *Queue {
	return &Queue{}
}

func getProtocol(env string) string {
	if env == "PROD" {
		return "amqps"
	}
	return "amqp"
}

func (q *Queue) dial(username, password, serverUrl string) {

	uri := fmt.Sprintf("%s://%s:%s@%s", getProtocol(os.Getenv("ENV")), username, password, serverUrl)
	fmt.Println(uri)
	conn, err := amqp.Dial(uri)
	if err != nil {
		log.Panicf("%s: %s", "Failed to connect to RabbitMQ", err)
	}
	q.conn = conn

	ch, err := q.conn.Channel()
	q.ch = ch
}

func (q *Queue) Setup(username, password, serverUrl string) error {

	// defer conn.Close()
	q.dial(username, password, serverUrl)

	logQ, err := q.ch.QueueDeclare(
		"LogMessages", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	q.LogQ = logQ

	// err check
	if err != nil {
		log.Panicf("%s: %s", "Failed to connect to RabbitMQ", err)
	}
	eventQ, err := q.ch.QueueDeclare(
		"EventMessages", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)

	// err check
	if err != nil {
		log.Panicf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

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

func (q *Queue) Consume(nameq string, callback callbackFunc) {
	// Consumir a fila
	msgs, err := q.ch.Consume(
		nameq, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		log.Fatalf("Falha ao consumir a fila: %s", err)
	}
	// Loop para consumir as mensagens
	for msg := range msgs {
		callback(msg)
	}
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
