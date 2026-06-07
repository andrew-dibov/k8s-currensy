package servers

import (
	"context"
	"conversion/internal/clients"
	"conversion/internal/producers"
	pb "conversion/internal/protos/conversion"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type ConversionServer struct {
	pb.UnimplementedConversionServiceServer
	currency *clients.CurrencyClient
	redis    *clients.RedisClient
	logger   *logrus.Logger
	producer *producers.KafkaProducer
}

func NewConversionServer(currencyClient *clients.CurrencyClient, redisClient *clients.RedisClient, kafkaProducer *producers.KafkaProducer, logger *logrus.Logger) *ConversionServer {
	return &ConversionServer{
		currency: currencyClient,
		redis:    redisClient,
		producer: kafkaProducer,
		logger:   logger,
	}
}

/* --- --- --- */

func (conversion *ConversionServer) Convert(ctx context.Context, req *pb.ConvertRequest) (*pb.ConvertResponse, error) {
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

	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	if req.Amount > 1e12 {
		return nil, fmt.Errorf("amount is too large")
	}

	rate, err := conversion.getRate(
		ctx,
		req.FromCurrency,
		req.ToCurrency,
	)

	if err != nil {
		return nil, fmt.Errorf("conversion failed : %w", err)
	}

	result := req.Amount * rate

	conversion.logger.WithFields(logrus.Fields{
		"from":   req.FromCurrency,
		"to":     req.ToCurrency,
		"amount": req.Amount,
		"result": result,
		"rate":   rate,
	}).Info("conversion")

	res := &pb.ConvertResponse{
		Result:       result,
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         rate,
		Amount:       req.Amount,
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := conversion.producer.SendConversion(ctx, res); err != nil {
			conversion.logger.WithError(err).Error("failed to send conversion event to Kafka")
		}
	}()

	return res, nil
}

func (conversion *ConversionServer) getRate(ctx context.Context, fromCurrency string, toCurrency string) (float64, error) {
	if fromCurrency == toCurrency {
		return 1.0, nil
	}

	if rate, found, err := conversion.redis.GetRate(ctx, fromCurrency, toCurrency); err == nil {
		if found {
			return rate, nil
		}
	} else {
		conversion.logger.WithError(err).Warn("failed to get from cache, fallback to currency")
	}

	data, err := conversion.currency.GetRate(ctx, fromCurrency, toCurrency)
	if err != nil {
		return 0, err
	}
	if data.Rate <= 0 {
		return 0, fmt.Errorf("received invalid rate : %f", data.Rate)
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := conversion.redis.SetRate(ctx, fromCurrency, toCurrency, data.Rate); err != nil {
			conversion.logger.WithError(err).Error("failed to cache rate")
		}
	}()

	return data.Rate, nil
}
