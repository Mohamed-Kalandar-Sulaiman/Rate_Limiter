package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// RateLimiterRepository defines methods for interacting with Redis for rate limiting.
type RateLimiterRepository struct {
	client *redis.Client
}

// NewRateLimiterRepository creates a new RateLimiterRepository instance.
func NewRateLimiterRepository(client *redis.Client) *RateLimiterRepository {
	return &RateLimiterRepository{
		client: client,
	}
}

// SetRequestCount sets a value for a given key and applies a TTL.
func (r *RateLimiterRepository) Set(key string, value int, ttl time.Duration) (string, error) {
    // Set the value of the key in Redis with a TTL.
    result, err := r.client.Set(context.Background(), key, value, ttl).Result()
    if err != nil {
        return "", fmt.Errorf("failed to set request count with TTL in Redis: %v", err)
    }
    return result, nil
}


// IncrementRequestCount increments the request count for a given key.
func (r *RateLimiterRepository) Increment(key string) (int64, error) {
	// Increment the count by 1 and return the new count.
	count, err := r.client.Incr(context.Background(), key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment request count in Redis: %v", err)
	}
	return count, nil
}

// GetTTL retrieves the TTL (Time to Live) for a given key.
func (r *RateLimiterRepository) GetTTL(key string) (time.Duration, error) {
	// Get the TTL of the key in Redis
	ttl, err := r.client.TTL(context.Background(), key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get TTL for key from Redis: %v", err)
	}

	return ttl, nil
}


// SetTTL sets the TTL (Time to Live) for a given key in Redis.
func (r *RateLimiterRepository) SetTTL(key string, ttl time.Duration) error {
	// Set the TTL for the key
	err := r.client.Expire(context.Background(), key, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set TTL for key in Redis: %v", err)
	}

	return nil
}


func (r *RateLimiterRepository) Get(key string) (int64, error) {
	// Get the current value of the request count from Redis.
	count, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		// If the key does not exist, return 0
		return -1, nil
	} else if err != nil {
		// If an error occurs other than Nil, return an error
		return 0, fmt.Errorf("failed to get request count from Redis: %v", err)
	}

	// Convert the string count to int64
	countInt, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert count to int64: %v", err)
	}

	return countInt, nil
}

// ClearKey removes the given key from Redis.
func (r *RateLimiterRepository) ClearKey(key string) error {
	// Delete the key from Redis to clear its value.
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

