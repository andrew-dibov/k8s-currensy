package clients

import (
	pb "api-gateway/internal/protos/conversion"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type ConversionClient struct {
	pbuf pb.ConversionServiceClient
	conn *grpc.ClientConn
}

func NewConversionClient(url string) (*ConversionClient, error) {
	conn, err := grpc.NewClient(url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{
		"loadBalancingPolicy": "round_robin",
    "methodConfig": [{
      "name": [{"service": "conversion.ConversionService"}],
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

	return &ConversionClient{pbuf: pb.NewConversionServiceClient(conn), conn: conn}, nil
}

func (c *ConversionClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

/* --- --- --- */

func (c *ConversionClient) Convert(ctx context.Context, fromCurrency string, toCurrency string, amount float64) (*pb.ConvertResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.pbuf.Convert(ctx, &pb.ConvertRequest{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Amount:       amount,
	})
}
