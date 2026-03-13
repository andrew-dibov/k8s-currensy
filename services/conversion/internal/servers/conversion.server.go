package servers

import (
	"context"
	"conversion/internal/clients"
	pb "conversion/internal/protos/conversion"
	"fmt"

	"github.com/sirupsen/logrus"
)

type ConversionServer struct {
	pb.UnimplementedConversionServiceServer
	redis *clients.RedisClient
	curr  *clients.CurrencyClient
	log   *logrus.Logger
}

func NewConversionServer(curr *clients.CurrencyClient, redis *clients.RedisClient, log *logrus.Logger) *ConversionServer {
	return &ConversionServer{
		redis: redis,
		curr:  curr,
		log:   log,
	}
}

/* --- --- --- */

func (cs *ConversionServer) Convert(ctx context.Context, req *pb.ConvertRequest) (*pb.ConvertResponse, error) {
	if req.FromCurrency == "" || req.ToCurrency == "" || req.Amount <= 0 {
		return nil, fmt.Errorf("invalid params")
	}

	rate, err := cs.getRate(ctx, req.FromCurrency, req.ToCurrency)
	if err != nil {
		return nil, fmt.Errorf("conversion failed : %w", err)
	}

	result := req.Amount * rate

	cs.log.WithFields(logrus.Fields{
		"from":   req.FromCurrency,
		"to":     req.ToCurrency,
		"amount": req.Amount,
		"result": result,
		"rate":   rate,
	}).Info("convertions performed")

	return &pb.ConvertResponse{
		Result:       result,
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         rate,
		Amount:       req.Amount,
	}, nil
}

func (cs *ConversionServer) getRate(ctx context.Context, fromCurrency string, toCurrency string) (float64, error) {
	if fromCurrency == toCurrency {
		return 1.0, nil
	}

	if rate, found, err := cs.redis.GetRate(ctx, fromCurrency, toCurrency); err == nil && found {
		return rate, nil
	}

	data, err := cs.curr.GetRate(ctx, fromCurrency, toCurrency)
	if err != nil {
		return 0, err
	}

	go func() {
		if err := cs.redis.SetRate(context.Background(), fromCurrency, toCurrency, data.Rate); err != nil {
			cs.log.WithError(err).Error("failed to cache rate")
		}
	}()

	return data.Rate, nil
}
