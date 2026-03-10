package clients

import (
	pb "api-gateway/internal/protos/rates"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type RatesClient struct {
	pbuf pb.RatesServiceClient
	conn *grpc.ClientConn
}

/* --- --- --- */

func NewRatesClient(url string) (*RatesClient, error) {
	conn, err := grpc.NewClient(url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{
		"loadBalancingPolicy": "round_robin",
    "methodConfig": [{
      "name": [{"service": "rates.RatesService"}],
      "retryPolicy": {
        "maxAttempts": 3,
        "initialBackoff": "0.1s",
        "maxBackoff": "1s",
        "backoffMultiplier": 2,
        "retryableStatusCodes": ["UNAVAILABLE"]
      }
    }]
	}`))

	if err != nil {
		return nil, fmt.Errorf("failed to establish connection : %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		state := conn.GetState()

		if state == connectivity.Ready {
			break
		}

		if state == connectivity.Shutdown {
			return nil, fmt.Errorf("connection shutdown")
		}

		if !conn.WaitForStateChange(ctx, state) {
			if ctx.Err() == context.DeadlineExceeded {
				return nil, fmt.Errorf("connection timeout")
			}
			return nil, fmt.Errorf("connection failed : %w", ctx.Err())
		}
	}

	return &RatesClient{pbuf: pb.NewRatesServiceClient(conn), conn: conn}, nil
}

func (client *RatesClient) Close() error {
	if client.conn != nil {
		return client.conn.Close()
	}
	return nil
}

/* --- --- --- */

func (client *RatesClient) GetRate(ctx context.Context, fromCurrency string, toCurrency string) (*pb.GetRateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return client.pbuf.GetRate(ctx, &pb.GetRateRequest{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
	})
}

func (client *RatesClient) GetAllRates(ctx context.Context, baseCurrency string) (*pb.GetAllRatesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return client.pbuf.GetAllRates(ctx, &pb.GetAllRatesRequest{
		BaseCurrency: baseCurrency,
	})
}

func (client *RatesClient) Convert(ctx context.Context, fromCurrency string, toCurrency string, amount float64) (*pb.ConvertResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return client.pbuf.Convert(ctx, &pb.ConvertRequest{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Amount:       amount,
	})
}
