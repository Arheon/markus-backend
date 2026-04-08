package kafka

import (
	"context"
	"fmt"

	"github.com/Arheon/markus-backend/pkg/outbox"
	"github.com/IBM/sarama"
)

type Broker struct {
	producer sarama.SyncProducer
	topics   map[string]string
}

type BrokerOption func(*Broker)

func WithTopics(topics map[string]string) BrokerOption {
	return func(b *Broker) {
		b.topics = topics
	}
}

func NewBroker(producer sarama.SyncProducer, opts ...BrokerOption) *Broker {
	b := &Broker{
		producer: producer,
		topics:   make(map[string]string),
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func (k *Broker) Publish(ctx context.Context, msg outbox.Message) error {
	topic, ok := k.topics[msg.Type]
	if !ok {
		return fmt.Errorf("no topic for type %s", msg.Type)
	}

	_, _, err := k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msg.Payload),
	})

	return err
}
