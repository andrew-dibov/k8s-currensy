package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	ttl    time.Duration
	client *redis.Client
}

func NewRedisClient(addr string, pass string, db int) (*RedisClient, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := rc.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed : %w", err)
	}

	return &RedisClient{
		ttl:    5 * time.Minute,
		client: rc,
	}, nil
}

func (rc *RedisClient) Close() error {
	return rc.client.Close()
}

/* --- --- --- */

func (rc *RedisClient) GetRate(ctx context.Context, fromCurrency string, toCurrency string) (float64, bool, error) {
	key := fmt.Sprintf("rate:%s:%s", fromCurrency, toCurrency)

	val, err := rc.client.Get(ctx, key).Result()

	if err == redis.Nil {
		return 0, false, nil
	}

	if err != nil {
		return 0, false, err
	}

	var rate float64
	if err := json.Unmarshal([]byte(val), &rate); err != nil {
		return 0, false, err
	}

	return rate, true, nil
}

func (rc *RedisClient) SetRate(ctx context.Context, fromCurrency string, toCurrency string, rate float64) error {
	key := fmt.Sprintf("rate:%s:%s", fromCurrency, toCurrency)

	val, err := json.Marshal(rate)
	if err != nil {
		return err
	}

	return rc.client.Set(ctx, key, val, rc.ttl).Err()
}
