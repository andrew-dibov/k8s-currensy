package clients

import (
	"context"
	pb "conversion/internal/protos/currency"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type CurrencyClient struct {
	client  pb.CurrencyServiceClient
	conn    *grpc.ClientConn
	timeout time.Duration
}

func NewCurrencyClient(currencyURL string, timeout time.Duration) (*CurrencyClient, error) {
	conn, err := grpc.NewClient(currencyURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: 2 * time.Second,
		}),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    10 * time.Second,
			Timeout: 1 * time.Second,
		}),
		grpc.WithDefaultServiceConfig(`{
			"loadBalancingPolicy": "round_robin",
			"methodConfig": [{
				"name": [{"service": "currency.CurrencyService"}],
				"retryPolicy": {
					"maxAttempts": 3,
					"maxBackoff": "1s",
					"backoffMultiplier": 2,
					"initialBackoff": "0.1s",
					"retryableStatusCodes": ["UNAVAILABLE"]
				}
			}]
		}`),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(4*1024*1024),
			grpc.MaxCallSendMsgSize(4*1024*1024),
		))

	if err != nil {
		return nil, fmt.Errorf("client init failed : %w", err)
	}

	conn.Connect()

	return &CurrencyClient{
		client:  pb.NewCurrencyServiceClient(conn),
		conn:    conn,
		timeout: timeout,
	}, nil
}

func (currency *CurrencyClient) Close() error {
	if currency.conn != nil {
		return currency.conn.Close()
	}
	return nil
}

/* --- --- --- */

func (currency *CurrencyClient) GetRate(ctx context.Context, fromCurrency string, toCurrency string) (*pb.GetRateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, currency.timeout)
	defer cancel()

	return currency.client.GetRate(ctx, &pb.GetRateRequest{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
	})
}

func (currency *CurrencyClient) GetAllRates(ctx context.Context, baseCurrency string) (*pb.GetAllRatesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, currency.timeout)
	defer cancel()

	return currency.client.GetAllRates(ctx, &pb.GetAllRatesRequest{
		BaseCurrency: baseCurrency,
	})
}
