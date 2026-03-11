package main

import (
	"context"
	"currency/internal/configs"
	"currency/internal/fetchers"
	"currency/internal/repositories"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"currency/internal/protos/currency"

	_ "github.com/lib/pq"

	"google.golang.org/grpc"
)

type Server struct {
	currency.UnimplementedCurrencyServiceServer
	repo    *repositories.Postgres
	fetcher *fetchers.ExchangeRateFetcher
}

func main() {
	cfg := configs.Load()

	/* --- */

	db, err := sql.Open("postgres", cfg.Postgres)
	if err != nil {
		log.Fatal("failed to connect to postgres : ", err)
	}
	defer db.Close()

	/* --- */

	repo := repositories.NewPostgres(db)
	fetcher := fetchers.NewExchangeRateFetcher(cfg.ExternalAPIURL, cfg.ExternalAPIToken)

	/* --- */

	go startUpdates(repo, fetcher)

	/* --- */

	srv := &Server{
		repo:    repo,
		fetcher: fetcher,
	}

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatal("failed to listen")
	}

	/* GRPC SERVER */

	grpcServer := grpc.NewServer()
	currency.RegisterCurrencyServiceServer(grpcServer, srv)

	/* RUN SERVER */

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed", err)
	}

}

func startUpdates(repo *repositories.Postgres, fetcher *fetchers.ExchangeRateFetcher) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	updateRates(repo, fetcher)

	for range ticker.C {
		updateRates(repo, fetcher)
	}
}

func updateRates(repo *repositories.Postgres, fetcher *fetchers.ExchangeRateFetcher) {
	rates, err := fetcher.FetchRates("USD")
	if err != nil {
		log.Printf("fail")
		return
	}

	if err := repo.UpdateRates(context.Background(), "USD", rates); err != nil {
		log.Printf("fail")
		return
	}
}

/* --- --- --- */

func (srv *Server) GetRate(ctx context.Context, req *currency.GetRateRequest) (*currency.GetRateResponse, error) {
	rate, err := srv.repo.GetRate(ctx, req.FromCurrency, req.ToCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to get rate : %v", err)
	}

	return &currency.GetRateResponse{
		Rate:         rate,
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
	}, nil
}

func (srv *Server) GetAllRates(ctx context.Context, req *currency.GetAllRatesRequest) (*currency.GetAllRatesResponse, error) {
	rates, err := srv.repo.GetAllRates(ctx, req.BaseCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed %v", err)
	}

	return &currency.GetAllRatesResponse{
		BaseCurrency: req.BaseCurrency,
		Rates:        rates,
	}, nil
}
