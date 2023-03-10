package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/badico-cloud-hub/log-driver/infra"
	"github.com/badico-cloud-hub/log-driver/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getEnvOrDefault(env, deflt string) string {
	v := os.Getenv(env)
	if v == "" {
		return deflt
	}
	return v
}

func main() {
	// Conectar ao servidor RabbitMQ
	infinite := make(chan bool)

	client, err := mongo.NewClient(
		options.Client().ApplyURI(
			getEnvOrDefault("MONGODB_URL", "mongodb://localhost:27017"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	logsCol := client.Database("prod").Collection("logs")
	eventsCol := client.Database("prod").Collection("events")

	if err != nil {
		log.Fatal(err)
	}

	q := infra.NewQueue()
	q.Setup(
		getEnvOrDefault("RABBITMQ_USERNAME", "user"),
		getEnvOrDefault("RABBITMQ_PASSWORD", "password"),
		getEnvOrDefault("RABBITMQ_URL", "localhost:5672"),
	)

	go q.Consume("EventMessages", func(msg amqp.Delivery) {
		var eventMessage logger.LogEventMessage
		json.Unmarshal(msg.Body, &eventMessage)
		_, err = eventsCol.InsertOne(context.TODO(), eventMessage)
	})

	go q.Consume("LogMessages", func(msg amqp.Delivery) {
		var logMessage logger.LogMessage
		json.Unmarshal(msg.Body, &logMessage)
		_, err = logsCol.InsertOne(context.TODO(), logMessage)
	})
	<-infinite
}
