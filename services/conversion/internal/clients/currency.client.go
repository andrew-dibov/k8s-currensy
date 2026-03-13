package clients

import (
	"context"
	pb "conversion/internal/protos/currency"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type CurrencyClient struct {
	pbuf pb.CurrencyServiceClient
	conn *grpc.ClientConn
}

func NewCurrencyClient(url string) (*CurrencyClient, error) {
	conn, err := grpc.NewClient(url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{
		"loadBalancingPolicy": "round_robin",
    "methodConfig": [{
      "name": [{"service": "currency.CurrencyService"}],
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

	conn.Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

	return &CurrencyClient{pbuf: pb.NewCurrencyServiceClient(conn), conn: conn}, nil
}

func (c *CurrencyClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

/* --- --- --- */

func (c *CurrencyClient) GetRate(ctx context.Context, fromCurrency string, toCurrency string) (*pb.GetRateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.pbuf.GetRate(ctx, &pb.GetRateRequest{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
	})
}
