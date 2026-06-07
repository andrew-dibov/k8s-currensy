package clients

import (
	pb "api-gateway/internal/protos/conversion"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type ConversionClient struct {
	client  pb.ConversionServiceClient
	conn    *grpc.ClientConn
	timeout time.Duration
}

func NewConversionClient(conversionURL string, timeout time.Duration) (*ConversionClient, error) {
	conn, err := grpc.NewClient(conversionURL,
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
				"name": [{"service": "conversion.ConversionService"}],
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

	return &ConversionClient{
		client:  pb.NewConversionServiceClient(conn),
		conn:    conn,
		timeout: timeout,
	}, nil
}

func (conversion *ConversionClient) Close() error {
	if conversion.conn != nil {
		return conversion.conn.Close()
	}
	return nil
}

/* --- --- --- */

func (conversion *ConversionClient) Convert(ctx context.Context, fromCurrency string, toCurrency string, amount float64) (*pb.ConvertResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, conversion.timeout)
	defer cancel()

	return conversion.client.Convert(ctx, &pb.ConvertRequest{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Amount:       amount,
	})
}
