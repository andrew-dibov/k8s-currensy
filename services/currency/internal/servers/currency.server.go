package servers

import (
	"context"
	proto "currency/internal/protos/currency"
	"currency/internal/repos"
	"fmt"
)

type CurrencyServer struct {
	proto.UnimplementedCurrencyServiceServer
	rp *repos.Postgres
}

func NewCurrencyServer(rp *repos.Postgres) *CurrencyServer {
	return &CurrencyServer{
		rp: rp,
	}
}

/* --- --- --- */

func (cs *CurrencyServer) GetRate(ctx context.Context, req *proto.GetRateRequest) (*proto.GetRateResponse, error) {
	rate, err := cs.rp.GetRate(ctx, req.FromCurrency, req.ToCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to get rate : %v", err)
	}

	return &proto.GetRateResponse{
		Rate:         rate,
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
	}, nil
}

func (cs *CurrencyServer) GetAllRates(ctx context.Context, req *proto.GetAllRatesRequest) (*proto.GetAllRatesResponse, error) {
	rates, err := cs.rp.GetAllRates(ctx, req.BaseCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to get rates : %v", err)
	}

	return &proto.GetAllRatesResponse{
		Rates:        rates,
		BaseCurrency: req.BaseCurrency,
	}, nil
}
