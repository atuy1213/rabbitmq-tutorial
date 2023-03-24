package main

import (
	"fmt"
	"log"

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

	msgs, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err, "Failed to consume a queue")
	}

	var forever chan struct{}
	go func() {
		for d := range msgs {
			log.Printf("Recieved a message: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for message. To exit press ctrl+c")
	<-forever
}
