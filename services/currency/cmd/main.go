package main

import (
	"context"
	"currency/internal/clients"
	"currency/internal/configs"
	"currency/internal/repos"
	"currency/internal/servers"
	"database/sql"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "currency/internal/protos/currency"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	appConfig := configs.LoadAppConfig()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.WithFields(logrus.Fields{
		"app_port":       appConfig.AppPort,
		"external_url":   appConfig.ExternalURL,
		"external_token": appConfig.ExternalToken,
		"postgres_url":   appConfig.PostgresURL,
	}).Info("currency : app config loaded")

	/* --- --- --- */

	db, err := sql.Open("postgres", appConfig.PostgresURL)
	if err != nil {
		logger.WithError(err).Fatal("failed to open postgres")
	}
	defer db.Close()

	/* --- --- --- */

	postgresRepo := repos.NewPostgresRepo(db)
	exchangeClient := clients.NewExchangeClient(appConfig.ExternalURL, appConfig.ExternalToken)

	go startUpdates(postgresRepo, exchangeClient, logger)

	/* --- --- --- */

	srv := servers.NewCurrencyServer(postgresRepo, logger)
	grpcSrv := grpc.NewServer(grpc.MaxRecvMsgSize(4*1024*1024), grpc.MaxSendMsgSize(4*1024*1024),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    10 * time.Second,
			Timeout: 1 * time.Second,
		}))

	pb.RegisterCurrencyServiceServer(grpcSrv, srv)

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

/* --- --- --- */

func startUpdates(repo *repos.PostgresRepo, exchange *clients.ExchangeClient, logger *logrus.Logger) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	updateRates(repo, exchange, logger)
	for range ticker.C {
		updateRates(repo, exchange, logger)
	}
}

func updateRates(repo *repos.PostgresRepo, exchange *clients.ExchangeClient, logger *logrus.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rates, err := exchange.GetRates(ctx, "USD")
	if err != nil {
		logger.WithError(err).Error("update rates failed")
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	if err := repo.UpdateRates(ctx, "USD", rates); err != nil {
		logger.WithError(err).Error("update rates in postgres failed")
		return
	}
}
