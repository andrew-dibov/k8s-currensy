package clients

import (
	"context"
	"fmt"
	"strings"
	"time"

	rd "github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *rd.Client
	ttl    time.Duration
}

func NewRedisClient(redisDB int, redisURL string, redisPass string, ttl time.Duration) (*RedisClient, error) {
	client := rd.NewClient(&rd.Options{
		DB:       redisDB,
		Addr:     redisURL,
		Password: redisPass,

		MaxRetries:      3,
		MinRetryBackoff: 100 * time.Millisecond,
		MaxRetryBackoff: 2 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis failed to ping : %w", err)
	}

	return &RedisClient{
		client: client,
		ttl:    ttl,
	}, nil
}

func (redis *RedisClient) Close() error {
	return redis.client.Close()
}

/* --- --- --- */

func normalizeCurrency(currency string) string {
	return strings.ToUpper(strings.TrimSpace(currency))
}

func (redis *RedisClient) GetRate(ctx context.Context, fromCurrency string, toCurrency string) (float64, bool, error) {
	from := normalizeCurrency(fromCurrency)
	to := normalizeCurrency(toCurrency)

	key := fmt.Sprintf("rate:%s:%s", from, to)
	rate, err := redis.client.Get(ctx, key).Float64()

	if err == rd.Nil {
		return 0, false, nil
	}

	if err != nil {
		return 0, false, fmt.Errorf("failed to get rate : %w", err)
	}

	return rate, true, nil
}

func (redis *RedisClient) SetRate(ctx context.Context, fromCurrency string, toCurrency string, rate float64) error {
	from := normalizeCurrency(fromCurrency)
	to := normalizeCurrency(toCurrency)

	key := fmt.Sprintf("rate:%s:%s", from, to)
	return redis.client.Set(ctx, key, rate, redis.ttl).Err()
}
