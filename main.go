package main

import (
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	constants "kstreams/gs/constants"
	services "kstreams/gs/services"
)

const (
	KafkaServer = "localhost:9092"
	KafkaTopic  = "transaction"
)

func main() {

	fmt.Println(constants.TransactionMock[0])

	// Set up Kafka consumer and producer configurations
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": KafkaServer,
		"group.id":          "aml-group",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	// Subscribe to the input topic
	err = consumer.SubscribeTopics([]string{KafkaTopic}, nil)
	if err != nil {
		panic(err)
	}

	for {
		msg, err := consumer.ReadMessage(0)
		if err != nil {
			panic(err)
		}

		record := strings.Split(string(msg.Value), ",")

		transaction, err := services.ParseTransaction(record)
		if err != nil {
			fmt.Printf("Error parsing transaction: %v\n", err)
			continue
		}

		topic := KafkaTopic

		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(constants.TransactionMock[0]),
		}, nil)

		if err != nil {
			panic(err)
		}

		if transaction.CostPerItem*float64(transaction.NumberOfItemsPurchased) > 1000 {
			// Potential AML alert, send to alert topic
			alertMsg := fmt.Sprintf("Potential AML alert for TransactionId: %s", transaction.TransactionId)
			deliveryChan := make(chan kafka.Event)
			err := producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &alertMsg, Partition: kafka.PartitionAny},
				Value:          []byte(alertMsg),
			}, deliveryChan)
			e := <-deliveryChan
			m := e.(*kafka.Message)
			if m.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
			} else {
				fmt.Println(err)
				fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
					*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
		}
	}
}
