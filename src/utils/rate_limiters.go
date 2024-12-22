package utils

import (
	"fmt"
	"time"

	"github.com/Mohamed-Kalandar-Sulaiman/Rate_Limiter/src/repository"
)


func ConvertUnitToTTL(unit string) (time.Duration, error) {
	var ttl time.Duration
	switch unit {
	case "seconds":
		ttl = time.Second
	case "minutes":
		ttl = time.Minute
	case "hours":
		ttl = time.Hour
	default:
		return 0, fmt.Errorf("unsupported unit: %v", unit)
	}
	return ttl, nil
}




type RateLimiter interface {
	RateLimitFunction(key, unit string, requestPerUnit int) (bool, int, int, int64, int64, error)
}

type FixedWindowRateLimiter struct {
	repo *repository.RateLimiterRepository
}

func NewFixedWindowRateLimiter(repo *repository.RateLimiterRepository) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{repo: repo}
}


type SlidingWindowRateLimiter struct {
	repo *repository.RateLimiterRepository
}

func NewSlidingWindowRateLimiter(repo *repository.RateLimiterRepository) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{repo: repo}
}



func (r *FixedWindowRateLimiter) RateLimitFunction(key, unit string, requestPerUnit int) (bool, int, int, int64, int64, error) {
	  //   isAllowed , limit, remaining, reset , reset_after ,err 
	isAllowed   := false
	limit       := requestPerUnit
	remaining   := 0
	// ! Check if Key exists
	count  ,      _ := r.repo.Get(key)

	if   count == -1 {
		ttl, _ := ConvertUnitToTTL(unit)
		count, err := r.repo.Set(key, 1, ttl)
		if err != nil{
			fmt.Println("error occured while setting ", count)
		}
		isAllowed = true
		remaining = requestPerUnit-1
		
	} else if count >= int64(requestPerUnit) || count == 0 {
		// ! Limit has reached
		isAllowed   = false
		remaining   = 0

	}else{
		current_count , _ := r.repo.Increment(key)
		isAllowed = true
		remaining = requestPerUnit- int(current_count)
	}


		ttl  , _    := r.repo.GetTTL(key)
		reset_after := int64(ttl.Seconds())
		reset       := int64(time.Now().Unix()) + int64(ttl.Seconds())


	return isAllowed, limit ,remaining ,reset,reset_after,nil
}









func (r *SlidingWindowRateLimiter) RateLimitFunction(key, unit string, requestPerUnit int) (bool, int, int, int64, int64, error) {
	now := time.Now().Unix()
	ttl, err := ConvertUnitToTTL(unit)
	if err != nil {
		fmt.Println("Error:", err)
	} 
	startTime := int64(ttl.Seconds()) 
	windowStart := now - startTime
	fmt.Println(windowStart)

	err = r.repo.ZAdd(key, now)
	if err != nil {
		return false, 0, 0, 0, 0, fmt.Errorf("failed to add request to Redis: %v", err)
	}

	err = r.repo.ZRem(key, windowStart)
	if err != nil {
		return false, 0, 0, 0, 0, fmt.Errorf("failed to remove old requests from Redis: %v", err)
	}

	count, err := r.repo.ZGet(key, windowStart)
	if err != nil {
		return false, 0, 0, 0, 0, fmt.Errorf("failed to count requests in Redis: %v", err)
	}

	isAllowed := count < int64(requestPerUnit)
	remaining := int(int64(requestPerUnit) - count)

	resetAfter := int64(1 - (now - windowStart))
	reset := now + resetAfter

	return isAllowed, requestPerUnit, remaining, reset, resetAfter, nil
}




type RateLimiterFactory struct {
	repo *repository.RateLimiterRepository
}

func NewRateLimiterFactory(repo *repository.RateLimiterRepository) *RateLimiterFactory {
	return &RateLimiterFactory{repo: repo}
}


func (f *RateLimiterFactory) CreateRateLimiter(rateLimiterType string) (RateLimiter, error) {
	if rateLimiterType == "fixed_window" {
		return NewFixedWindowRateLimiter(f.repo), nil
	}
	if rateLimiterType == "sliding_window" {
		return NewSlidingWindowRateLimiter(f.repo), nil
	}
	return nil, fmt.Errorf("unsupported rate limiter type: %s", rateLimiterType)
}


