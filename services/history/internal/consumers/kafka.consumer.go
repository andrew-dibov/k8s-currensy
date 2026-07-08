package consumers

import (
	"context"
	"errors"
	"history/internal/repos"
	"io"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"

	pb "history/internal/protos/history"
)

type KafkaConsumer struct {
	reader *kafka.Reader
	repo   *repos.PostgresRepo
	logger *logrus.Logger
}

func NewKafkaConsumer(brokers string, topic string, groupID string, repo *repos.PostgresRepo, logger *logrus.Logger) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokers},
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: 0,
		StartOffset:    kafka.FirstOffset,
	})

	return &KafkaConsumer{
		reader: reader,
		repo:   repo,
		logger: logger,
	}
}

func (consumer *KafkaConsumer) Close() error {
	return consumer.reader.Close()
}

/* --- --- --- */

func (consumer *KafkaConsumer) Run(ctx context.Context) {
	for {
		msg, err := consumer.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				consumer.logger.Info("context canceled, exiting")
				return
			}

			if errors.Is(err, io.EOF) {
				time.Sleep(1 * time.Second)
				continue
			}

			consumer.logger.WithError(err).Error("failed to read message")
			time.Sleep(1 * time.Second)
			continue
		}

		var rec pb.ConversionRecord
		if err := protojson.Unmarshal(msg.Value, &rec); err != nil {
			consumer.logger.WithError(err).Error("failed to unmarshal message")
			continue
		}

		if err := consumer.repo.SaveConversion(ctx, &rec); err != nil {
			consumer.logger.WithError(err).Error("failed to save message")
			continue
		}

		if err := consumer.reader.CommitMessages(ctx, msg); err != nil {
			consumer.logger.WithError(err).Error("failed to commit offset")
		}

		consumer.logger.WithFields(logrus.Fields{
			"msg_id":     rec.Id,
			"msg_from":   rec.FromCurrency,
			"msg_to":     rec.ToCurrency,
			"msg_amount": rec.Amount,
			"msg_result": rec.Result,
		}).Info("message saved")
	}
}
