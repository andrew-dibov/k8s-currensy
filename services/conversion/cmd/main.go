package main

import (
	"conversion/internal/clients"
	"conversion/internal/configs"
	"conversion/internal/producers"
	"conversion/internal/servers"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "conversion/internal/protos/conversion"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	appConfig := configs.LoadAppConfig()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.WithFields(logrus.Fields{
		"app_port":     appConfig.AppPort,
		"currency_url": appConfig.CurrencyURL,

		"redis_db":   appConfig.RedisDB,
		"redis_url":  appConfig.RedisURL,
		"redis_pass": appConfig.RedisPass,

		"kafka_topic":   appConfig.KafkaTopic,
		"kafka_brokers": appConfig.KafkaBrokers,
	}).Info("conversion : app config loaded")

	/* --- --- --- */

	currencyClient, err := clients.NewCurrencyClient(appConfig.CurrencyURL, 5*time.Second)
	if err != nil {
		logger.WithError(err).Fatal("currency client failed")
	}
	defer currencyClient.Close()

	redisClient, err := clients.NewRedisClient(appConfig.RedisDB, appConfig.RedisURL, appConfig.RedisPass, 5*time.Minute)
	if err != nil {
		logger.WithError(err).Fatal("redis client failed")
	}
	defer redisClient.Close()

	kafkaProducer, err := producers.NewKafkaProducer(appConfig.KafkaBrokers, appConfig.KafkaTopic, logger)
	if err != nil {
		logger.WithError(err).Fatal("kafka client failed")
	}
	defer kafkaProducer.Close()

	/* --- --- --- */

	srv := servers.NewConversionServer(currencyClient, redisClient, kafkaProducer, logger)
	grpcSrv := grpc.NewServer(grpc.MaxRecvMsgSize(4*1024*1024), grpc.MaxSendMsgSize(4*1024*1024),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    10 * time.Second,
			Timeout: 1 * time.Second,
		}))

	pb.RegisterConversionServiceServer(grpcSrv, srv)

	lis, err := net.Listen("tcp", appConfig.AppPort)
	if err != nil {
		logger.WithError(err).Fatal("listener failed")
	}

	go func() {
		logger.WithField("port", appConfig.AppPort).Info("server starting")
		if err := grpcSrv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			logger.WithError(err).Fatal("server failed")
		}
	}()

	/* --- --- --- */

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("server shutting down")

	done := make(chan struct{})
	go func() {
		grpcSrv.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		logger.Info("server stopped")
	case <-time.After(30 * time.Second):
		logger.Warn("server forced to stop")
		grpcSrv.Stop()
	}
}
