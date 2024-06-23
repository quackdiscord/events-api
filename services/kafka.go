package services

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/quackdiscord/events-api/structs"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	log "github.com/sirupsen/logrus"
)

var Kafka *KafkaService
var KafkaReader *KafkaReaderService
var Events = make(map[string]*Event)

type Event struct {
	Name    string
	Handler func(string)
}

type KafkaService struct {
	writer *kafka.Writer
}

type KafkaReaderService struct {
	reader *kafka.Reader
}

func ConnectKafka(events []*structs.Event) {
	broker := os.Getenv("KAFKA_BROKER")
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	topic := os.Getenv("KAFKA_TOPIC")

	mechanism, err := scram.Mechanism(scram.SHA256, username, password)
	if err != nil {
		log.WithError(err).Fatal("Error creating Kafka SASL mechanism")
	}

	KafkaReader = &KafkaReaderService{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: "events-group-1",
			Dialer: &kafka.Dialer{
				SASLMechanism: mechanism,
				TLS:           &tls.Config{},
			},
		}),
	}

	// add each event from events/handler.go to the map
	for _, event := range events {
		Events[event.Name] = (*Event)(event)
		log.Info("Added event " + event.Name)
	}

	log.Info("Connected to Kafka")
}

func (k *KafkaService) Produce(ctx context.Context, key, value []byte) error {
	err := k.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
	if err != nil {
		log.WithError(err).Error("Error producing Kafka message")
		return err
	}
	return nil
}

func (kr *KafkaReaderService) Consume(ctx context.Context) {
	for {
		msg, err := kr.reader.ReadMessage(ctx)
		if err != nil {
			log.WithError(err).Error("Error reading Kafka message")
			break
		}
		key := string(msg.Key)
		if _, ok := Events[key]; ok {
			Events[key].Handler(string(msg.Value))
		} else {
			log.Infof("Consumed message: key=%s", key)
		}
	}
}

func DisconnectKafka() {
	if err := KafkaReader.reader.Close(); err != nil {
		log.WithError(err).Error("Error closing Kafka reader")
	}

	log.Info("Disconnected from Kafka")
}
