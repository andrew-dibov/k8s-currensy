package main

import (
	"context"
	"database/sql"
	"history/internal/configs"
	"history/internal/consumers"
	"history/internal/repos"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/sirupsen/logrus"
)

func main() {
	appConfig := configs.LoadAppConfig()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.WithFields(logrus.Fields{
		"app_port":      appConfig.AppPort,
		"postgres_url":  appConfig.PostgresURL,
		"kafka_brokers": appConfig.KafkaBrokers,
		"kafka_group":   appConfig.KafkaGroup,
		"kafka_topic":   appConfig.KafkaTopic,
	}).Info("history : app config loaded")

	/* --- --- --- */

	db, err := sql.Open("postgres", appConfig.PostgresURL)
	if err != nil {
		logger.WithError(err).Fatal("failed to open postgres")
	}
	defer db.Close()

	/* --- --- --- */

	repo := repos.NewPostgresRepo(db)

	kafkaConsumer := consumers.NewKafkaConsumer(appConfig.KafkaBrokers, appConfig.KafkaTopic, appConfig.KafkaGroup, repo, logger)
	defer kafkaConsumer.Close()

	logger.Info("starting service")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		logger.Info("shutting down")
		cancel()
	}()

	done := make(chan any)
	go func() {
		kafkaConsumer.Run(ctx)
		close(done)
	}()

	select {
	case <-done:
		logger.Info("consumer stopped")
	case <-time.After(30 * time.Second):
		logger.Warn("consumer forced to stop")
		if err := kafkaConsumer.Close(); err != nil {
			logger.WithError(err).Warn("error closing consumer")
		}
		<-done
	}

	logger.Info("service stopped")
}
