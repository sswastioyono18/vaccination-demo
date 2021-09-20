package main

import (
	"encoding/json"
	"github.com/sswastioyono18/vaccination-demo/internal/app/domain/resident"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	// Define RabbitMQ server URL.
	amqpServerURL := "amqp://guest:guest@localhost:5672"

	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	// Opening a channel to our RabbitMQ instance over
	// the connection we have already established.
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	q, err := channelRabbitMQ.QueueDeclare(
		"NewVaccineRegistration",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	err = channelRabbitMQ.QueueBind(
		q.Name, // queue name
		"",                          // routing key
		"ResidentVaccination", // exchange
		false,
		nil,
	)

	// Subscribing to QueueService1 for getting messages.
	messages, err := channelRabbitMQ.Consume(
		q.Name, // queue name
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no local
		false,           // no wait
		nil,             // arguments
	)
	if err != nil {
		log.Println(err)
	}

	// Build a welcome message.
	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages")

	// Make a channel to receive messages into infinite loop.
	forever := make(chan bool)

	go func() {
		for message := range messages {
			// For example, show received message in a console.
			log.Printf(" > Received message: %s\n", message.Body)
			var residentData resident.Resident
			err = json.Unmarshal(message.Body, &residentData)
			if err != nil {
				return
			}

			log.Println("NIK:", residentData.NIK)
		}
	}()

	<-forever
}
