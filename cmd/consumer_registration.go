package main

import (
	"encoding/json"
	"fmt"
	"github.com/sswastioyono18/vaccination-demo/config"
	residentDomain "github.com/sswastioyono18/vaccination-demo/internal/app/domain/resident"
	"github.com/sswastioyono18/vaccination-demo/internal/app/infra"
	zlog "github.com/sswastioyono18/vaccination-demo/internal/app/middleware"
	"go.uber.org/zap"
	"log"
	"time"
)

func main() {
	// Define RabbitMQ server URL.
	zlog.NewLogger("PROD")
	zlogger := zlog.Logger
	appConfig, err := config.NewConfig()
	messageQueueUri := fmt.Sprintf("amqp://%s:%s@%s", appConfig.MQ.User, appConfig.MQ.Pass, appConfig.MQ.Host)

	residentExchange, err := infra.NewBrokerExchange(appConfig.MQ.Resident.Exchanges.ResidentVaccination, appConfig.MQ.Resident.Queues.Registration, messageQueueUri)
	if err != nil {
		log.Fatal("error during init mq", err)
	}
	defer residentExchange.Channel.Close()

	queueArgs := make(map[string]interface{})
	queueArgs["x-queue-type"] = "quorum"

	q, err := residentExchange.Channel.QueueDeclare(
		residentExchange.QueueName, // name
		true,                       // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		queueArgs,                  // arguments
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
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // arguments
	)
	if err != nil {
		log.Println(err)
	}

	// Build a welcome message.
	zlogger.Info("Successfully connected to RabbitMQ")
	zlogger.Info("Waiting for messages")

	// Make a channel to receive messages into infinite loop.
	forever := make(chan bool)

	go func() {
		for message := range messages {
			// For example, show received message in a console.
			time.Sleep(10 * time.Second)
			zlogger.Info(" > Received message: \n", zap.String("Body: ", string(message.Body)))
			var residentData residentDomain.Resident
			err = json.Unmarshal(message.Body, &residentData)
			if err != nil {
				return
			}

			log.Println("NIK:", residentData.NIK)
			if err = message.Ack(false); err != nil {
				zlogger.Error("error", zap.Error(err))
			} else {
				zlogger.Info("acked message")
			}
		}
	}()

	<-forever
}
