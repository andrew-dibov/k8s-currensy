package servers

import (
	"context"
	pb "currency/internal/protos/currency"
	"currency/internal/repos"
	"fmt"

	"github.com/sirupsen/logrus"
)

type CurrencyServer struct {
	pb.UnimplementedCurrencyServiceServer
	postgres *repos.PostgresRepo
	logger   *logrus.Logger
}

func NewCurrencyServer(repo *repos.PostgresRepo, logger *logrus.Logger) *CurrencyServer {
	return &CurrencyServer{
		postgres: repo,
		logger:   logger,
	}
}

/* --- --- --- */

func (currency *CurrencyServer) GetRate(ctx context.Context, req *pb.GetRateRequest) (*pb.GetRateResponse, error) {
	if req.FromCurrency == "" || req.ToCurrency == "" {
		return nil, fmt.Errorf("fromCurrency or toCurrency is empty")
	}

	if len(req.FromCurrency) != 3 || len(req.ToCurrency) != 3 {
		return nil, fmt.Errorf("fromCurrency or toCurrency length is not 3 chars")
	}

	for _, ch := range req.FromCurrency {
		if ch < 'A' || ch > 'Z' {
			return nil, fmt.Errorf("fromCurrency has invalid chars")
		}
	}

	for _, ch := range req.ToCurrency {
		if ch < 'A' || ch > 'Z' {
			return nil, fmt.Errorf("toCurrency has invalid chars")
		}
	}

	rate, err := currency.postgres.GetRate(
		ctx,
		req.FromCurrency,
		req.ToCurrency,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get rate : %v", err)
	}

	return &pb.GetRateResponse{
		Rate:         rate,
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
	}, nil
}

func (currency *CurrencyServer) GetAllRates(ctx context.Context, req *pb.GetAllRatesRequest) (*pb.GetAllRatesResponse, error) {
	if req.BaseCurrency == "" {
		return nil, fmt.Errorf("baseCurrency is empty")
	}

	if len(req.BaseCurrency) != 3 {
		return nil, fmt.Errorf("baseCurrency length is not 3 chars")
	}

	for _, ch := range req.BaseCurrency {
		if ch < 'A' || ch > 'Z' {
			return nil, fmt.Errorf("baseCurrency has invalid chars")
		}
	}

	rates, err := currency.postgres.GetAllRates(ctx, "USD")
	if err != nil {
		currency.logger.WithError(err).Error("failed to get USD rates")
		return nil, fmt.Errorf("failed to get rates : %v", err)
	}

	if req.BaseCurrency == "USD" {
		if _, ok := rates["USD"]; !ok {
			rates["USD"] = 1.0
		}
		return &pb.GetAllRatesResponse{
			Rates:        rates,
			BaseCurrency: req.BaseCurrency,
		}, nil
	}

	bases, ok := rates[req.BaseCurrency]
	if !ok {
		currency.logger.WithField("baseCurrency", req.BaseCurrency).Warn("base currency not found in USD rates")
		return nil, fmt.Errorf("base currency not found : %s", req.BaseCurrency)
	}

	newRates := make(map[string]float64, len(rates))
	for code, rateUSD := range rates {
		newRates[code] = rateUSD / bases
	}

	return &pb.GetAllRatesResponse{
		Rates:        newRates,
		BaseCurrency: req.BaseCurrency,
	}, nil
}
