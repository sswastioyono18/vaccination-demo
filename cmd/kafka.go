package main

import (
	"context"
	"crypto/tls"
	"github.com/segmentio/kafka-go"

	//"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/segmentio/kafka-go/sasl/scram"
	"log"
	"time"
)

func main() {
	// get kafka writer using environment variables.
	kafkaURL := "sulky-02.srvs.cloudkafka.com:9094"
	topic := "q0u6jjmj-default"
	username := "q0u6jjmj"
	password := "RVpQjFgo5JOKczjJBTXjaVd51cEgOM-e"

	mechanism, err := scram.Mechanism(scram.SHA256, username, password)
	if err != nil {
		panic(err)
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
		TLS: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	conn, err := dialer.DialLeader(context.Background(), "tcp", kafkaURL, topic, 0)

	_, err = conn.WriteMessages(
		kafka.Message{
			Key:   []byte("Key-A"),
			Value: []byte("Hello World!"),
		},
		kafka.Message{
			Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		kafka.Message{
			Key:   []byte("Key-C"),
			Value: []byte("Two!"),
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	log.Println("finished")

	//config := &kafka.ConfigMap{
	//	"metadata.broker.list": kafkaURL,
	//	"security.protocol":    "SASL_SSL",
	//	"sasl.mechanisms":      "SCRAM-SHA-256",
	//	"sasl.username":        username,
	//	"sasl.password":       password,
	//	"group.id":             topic,
	//	"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"},
	//	//"debug":                           "generic,broker,security",
	//}
	////topic := os.Getenv("CLOUDKARAFKA_TOPIC_PREFIX") + ".test"
	//p, err := kafka.NewProducer(config)
	//if err != nil {
	//	fmt.Printf("Failed to create producer: %s\n", err)
	//	os.Exit(1)
	//}
	//fmt.Printf("Created Producer %v\n", p)
	//deliveryChan := make(chan kafka.Event)
	//
	//for i := 0; i < 10; i++ {
	//	value := fmt.Sprintf("[%d] Hello Go!", i+1)
	//	err = p.Produce(&kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, Value: []byte(value)}, deliveryChan)
	//	e := <-deliveryChan
	//	m := e.(*kafka.Message)
	//	if m.TopicPartition.Error != nil {
	//		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	//	} else {
	//		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
	//			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	//	}
	//}
	//close(deliveryChan)
}