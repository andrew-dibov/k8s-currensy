package main

import (
	"conversion/internal/clients"
	"conversion/internal/configs"
	"conversion/internal/servers"
	"net"

	pb "conversion/internal/protos/conversion"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	cfg := configs.Load()

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	currencyClient, err := clients.NewCurrencyClient(cfg.CurrencyService)
	if err != nil {
		log.WithError(err).Fatal("failed to create currency client")
	}
	defer currencyClient.Close()

	redisClient, err := clients.NewRedisClient(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDB)
	if err != nil {
		log.WithError(err).Fatal("failed to create redis client")
	}
	defer redisClient.Close()

	srv := grpc.NewServer()
	srvc := servers.NewConversionServer(currencyClient, redisClient, log)

	pb.RegisterConversionServiceServer(srv, srvc)

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.WithError(err).Fatal("failed to listen")
	}

	log.WithField("port", cfg.Port).Info("starting conversion service")
	if err := srv.Serve(lis); err != nil {
		log.WithError(err).Fatal("server failed")
	}
}
