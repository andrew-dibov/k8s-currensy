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

	/* --- --- --- */

	currencyClient, err := clients.NewCurrencyClient(cfg.CurrencyService)
	if err != nil {
		log.WithError(err).Fatal("currency client failed")
	}
	defer currencyClient.Close()

	redisClient, err := clients.NewRedisClient(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDB)
	if err != nil {
		log.WithError(err).Fatal("redis client failed")
	}
	defer redisClient.Close()

	/* --- --- --- */

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.WithError(err).Fatal("failed to announce listener")
	}

	srv := grpc.NewServer()
	srvc := servers.NewConversionServer(currencyClient, redisClient, log)

	pb.RegisterConversionServiceServer(srv, srvc)

	if err := srv.Serve(lis); err != nil {
		log.WithError(err).Fatal("failed to start conversion")
	}
}
