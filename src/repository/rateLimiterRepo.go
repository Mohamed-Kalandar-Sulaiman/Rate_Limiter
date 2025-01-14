package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RateLimiterRepository struct {
	client *redis.Client
}

func NewRateLimiterRepository(client *redis.Client) *RateLimiterRepository {
	return &RateLimiterRepository{
		client: client,
	}
}

func (r *RateLimiterRepository) Set(key string, value int, ttl time.Duration) (string, error) {
    result, err := r.client.Set(context.Background(), key, value, ttl).Result()
    if err != nil {
        return "", fmt.Errorf("failed to set request count with TTL in Redis: %v", err)
    }
    return result, nil
}


func (r *RateLimiterRepository) Increment(key string) (int64, error) {
	count, err := r.client.Incr(context.Background(), key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment request count in Redis: %v", err)
	}
	return count, nil
}

func (r *RateLimiterRepository) GetTTL(key string) (time.Duration, error) {
	ttl, err := r.client.TTL(context.Background(), key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get TTL for key from Redis: %v", err)
	}

	return ttl, nil
}


func (r *RateLimiterRepository) SetTTL(key string, ttl time.Duration) error {
	err := r.client.Expire(context.Background(), key, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set TTL for key in Redis: %v", err)
	}

	return nil
}


func (r *RateLimiterRepository) Get(key string) (int64, error) {
	count, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return -1, nil
	} else if err != nil {
		return 0, fmt.Errorf("failed to get request count from Redis: %v", err)
	}

	countInt, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert count to int64: %v", err)
	}

	return countInt, nil
}

func (r *RateLimiterRepository) ClearKey(key string) error {
	err := r.client.Del(context.Background(), key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete Redis key: %v", err)
	}
	return nil
}





 func (r *RateLimiterRepository) ZAdd(key string, value int64) error {
	_, err := r.client.ZAdd(context.Background(), key, &redis.Z{
		Score:  float64(value), 
		Member: float64(value),                
	}).Result()
	if err != nil {
		return fmt.Errorf("failed to add request timestamp to Redis: %v", err)
	}
	return nil
}

func (r *RateLimiterRepository) ZRem(key string, startTime int64) error {
	_, err := r.client.ZRemRangeByScore(context.Background(), key, "-inf", fmt.Sprintf("%d", startTime)).Result()
	if err != nil {
		return fmt.Errorf("failed to remove old request timestamps from Redis: %v", err)
	}
	return nil
}

func (r *RateLimiterRepository) ZGet(key string, startTime int64) (int64, error) {
	count, err := r.client.ZCount(context.Background(), key, fmt.Sprintf("%d", startTime), "+inf").Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get request count from Redis: %v", err)
	}
	return count, nil
}

