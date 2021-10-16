package main

import (
	"encoding/json"
	"fmt"
	"github.com/sswastioyono18/vaccination-demo/config"
	"github.com/sswastioyono18/vaccination-demo/internal/app/domain/resident"
	"github.com/sswastioyono18/vaccination-demo/internal/app/infra"
	"log"
)

func main() {
	// Define RabbitMQ server URL.
	appConfig, err := config.NewConfig()
	messageQueueUri := fmt.Sprintf("amqp://%s:%s@%s:%d",  appConfig.MQ.User,  appConfig.MQ.Pass,  appConfig.MQ.Host,  appConfig.MQ.Port)

	residentExchange,err  := infra.NewBrokerExchange(appConfig.MQ.Resident.Exchanges.ResidentVaccination, appConfig.MQ.Resident.Queues.Registration, messageQueueUri)
	if err != nil {
		log.Fatal("error during init mq", err)
	}
	defer residentExchange.Channel.Close()

	q, err := residentExchange.Channel.QueueDeclare(
		residentExchange.QueueName,    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	err = residentExchange.Channel.QueueBind(
		q.Name,                                     // queue name
		appConfig.MQ.Resident.Routing.Registration, // routing key
		residentExchange.ExchangeName,              // exchange
		false,
		nil,
	)

	// Subscribing to QueueService1 for getting messages.
	messages, err := residentExchange.Channel.Consume(
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
