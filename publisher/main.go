package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/atuy1213/rabbitmq-tutorial/pkg/repository"
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

	repository := repository.NewRepository()
	tasks := repository.GetTasks()

	if err := channel.Tx(); err != nil {
		log.Fatal(err, "Failed to start transaction")
	}

	for _, v := range tasks {

		if err := v.DoSomething(); err != nil {
			if err := channel.TxRollback(); err != nil {
				log.Fatal(err, "Failed to rollback transaction")
			}
			log.Fatal(err, "Failed to do something")
		}

		if err := channel.PublishWithContext(
			ctx,
			"",
			queue.Name,
			false,
			false,
			rabbitmq.Publishing{
				ContentType: "text/plain",
				Body:        []byte(v.Name),
			},
		); err != nil {
			if err := channel.TxRollback(); err != nil {
				log.Fatal(err, "Failed to rollback transaction")
			}
			log.Fatal(err, "Failed to publish messaige")
		}

		log.Printf("[x] Sending %s\n", v.Name)
	}

	if err := channel.TxCommit(); err != nil {
		log.Fatal(err, "Failed to commit transaction")
	}
}
