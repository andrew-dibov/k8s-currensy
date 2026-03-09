package clients

import (
	pb "api-gateway/internal/protos/rates"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type RatesClient struct {
	pbuf pb.RatesServiceClient
	conn *grpc.ClientConn
}

func NewRatesClient(url string) (*RatesClient, error) {
	conn, err := grpc.NewClient(url, grpc.WithDefaultServiceConfig(`{
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

func (client *RatesClient) GetRate(ctx context.Context, baseCurr string) (*pb.GetRateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return client.pbuf.GetRate(ctx, &pb.GetRateRequest{
		FromCurrency: "",
		ToCurrency:   "",
	})
}
