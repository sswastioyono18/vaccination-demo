package infra

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type MessageBroker struct {
	ExchangeName string
	QueueName string
	Channel *amqp.Channel
}

func (r MessageBroker) Publish(routingKey string, message []byte) (err error) {
	// Create a message to publish.
	messageBody := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType: "application/json",
		Body:        message,
	}

	// Attempt to publish a message to the routingKey.
	err = r.Channel.Publish(
		r.ExchangeName, // exchange
		routingKey,     // routingKey name
		false,          // mandatory
		false,          // immediate
		messageBody,    // message to publish
	)

	if err != nil {
		log.Fatal("error publish", err)
	}

	fmt.Println("Sending ", string(message))

	return
}

type IResidentExchange interface {
	Publish(queue string, message []byte) error
}

func NewBrokerExchange(exchange, queue , URIConnection string) (MessageBroker, error) {
	amqConnection, err := amqp.Dial(URIConnection)
	if err != nil {
		panic(err)
	}

	channelRabbitMQ, err := amqConnection.Channel()
	if err != nil {
		panic(err)
	}

	err = channelRabbitMQ.ExchangeDeclare(
		exchange, // name
		"direct",                    // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		panic(err)
	}

	return MessageBroker{
		exchange,
		queue,
		channelRabbitMQ,
	}, nil
}

