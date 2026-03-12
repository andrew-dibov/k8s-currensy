package main

import (
	"context"
	"currency/internal/clients"
	"currency/internal/configs"
	proto "currency/internal/protos/currency"
	"currency/internal/repos"
	"currency/internal/servers"
	"database/sql"
	"net"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	cfg := configs.Load()

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	/* --- --- --- */

	db, err := sql.Open("postgres", cfg.PostgresURL)
	if err != nil {
		log.WithError(err).Fatal("failed to open db")
	}
	defer db.Close()

	rp := repos.NewPostgres(db)
	ec := clients.NewExchangeClient(cfg.API, cfg.APIToken)

	go startUpdates(rp, ec, log)

	/* --- --- --- */

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.WithError(err).Fatal("failed to announce listener")
	}

	currency := servers.NewCurrencyServer(rp)
	server := grpc.NewServer()

	proto.RegisterCurrencyServiceServer(server, currency)

	if err := server.Serve(lis); err != nil {
		log.WithError(err).Error("failed to start currency server")
	}
}

func startUpdates(rp *repos.Postgres, ec *clients.ExchangeClient, log *logrus.Logger) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	updateRates(rp, ec, log)

	for range ticker.C {
		updateRates(rp, ec, log)
	}
}

func updateRates(rp *repos.Postgres, ec *clients.ExchangeClient, log *logrus.Logger) {
	rates, err := ec.GetRates("USD")
	if err != nil {
		log.WithError(err).Error("failed to get rates")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	if err := rp.UpdateRates(ctx, "USD", rates); err != nil {
		log.WithError(err).Error("failed to update rates")
		return
	}
}
