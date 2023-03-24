package main

import (
	"context"
	"fmt"
	"log"
	"time"

	rabbitmq "github.com/rabbitmq/amqp091-go"
)

const (
	RABBITMQ_USER     = "rabbitmq"
	RABBITMQ_PASSWORD = "password"
)

func main() {
	conn, err := rabbitmq.Dial(fmt.Sprintf("amqp://%s:%s@localhost:5672/", RABBITMQ_USER, RABBITMQ_PASSWORD))
	if err != nil {
		log.Fatal(err, "Failed to connect RabbitMQ")
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err, "Failed to open channel")
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err, "Failed to declare a queue")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World"
	if err := channel.PublishWithContext(
		ctx,
		"",
		queue.Name,
		false,
		false,
		rabbitmq.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	); err != nil {
		log.Fatal(err, "Failed to publish messaige")
	}

	log.Printf("[x] Sent %s\n", body)
}
