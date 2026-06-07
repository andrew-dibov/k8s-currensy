package producers

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"

	pb "conversion/internal/protos/conversion"
)

type KafkaProducer struct {
	writer *kafka.Writer
	logger *logrus.Logger
}

func NewKafkaProducer(brokers string, topic string, logger *logrus.Logger) (*KafkaProducer, error) {
	if err := ensureTopic(brokers, topic); err != nil {
		return nil, fmt.Errorf("failed to ensure topic : %w", err)
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		Async:        false,
	}

	return &KafkaProducer{
		writer: writer,
		logger: logger,
	}, nil
}

func (producer *KafkaProducer) Close() error {
	return producer.writer.Close()
}

func ensureTopic(brokers string, topic string) error {
	conn, err := kafka.Dial("tcp", brokers)
	if err != nil {
		return fmt.Errorf("failed to dial kafka : %w", err)
	}
	defer conn.Close()

	parts, err := conn.ReadPartitions(topic)
	if err == nil && len(parts) > 0 {
		return nil
	}

	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})

	if err != nil {
		return fmt.Errorf("failed to create topic : %w", err)
	}

	return nil
}

/* --- --- --- */

func (producer *KafkaProducer) SendConversion(ctx context.Context, record *pb.ConvertResponse) error {
	data, err := protojson.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record : %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(record.FromCurrency + record.ToCurrency),
		Value: data,
	}

	err = producer.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to write message : %w", err)
	}

	return nil
}
