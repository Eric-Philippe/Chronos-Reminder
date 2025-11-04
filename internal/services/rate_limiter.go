package services

import (
	"context"
	"fmt"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/database"
)

// RateLimiterService handles per-session rate limiting using Redis
type RateLimiterService struct {
	requestsPerWindow int
	windowDuration    time.Duration
}

// NewRateLimiterService creates a new rate limiter service
func NewRateLimiterService(requestsPerWindow int, windowSeconds int) *RateLimiterService {
	return &RateLimiterService{
		requestsPerWindow: requestsPerWindow,
		windowDuration:    time.Duration(windowSeconds) * time.Second,
	}
}

// RequestsPerWindow returns the configured requests per window limit
func (rl *RateLimiterService) RequestsPerWindow() int {
	return rl.requestsPerWindow
}

// WindowDuration returns the configured window duration
func (rl *RateLimiterService) WindowDuration() time.Duration {
	return rl.windowDuration
}

// RateLimitKey generates a cache key for rate limiting a session
func (rl *RateLimiterService) RateLimitKey(sessionID string) string {
	return fmt.Sprintf("rate_limit:%s", sessionID)
}

// RequestCounter holds the count of requests in the current window
type RequestCounter struct {
	Count     int       `json:"count"`
	ResetTime time.Time `json:"reset_time"`
}

// CheckRateLimit checks if a request should be allowed based on the session's rate limit
// Returns (allowed bool, remainingRequests int, resetTime time.Time, error)
func (rl *RateLimiterService) CheckRateLimit(ctx context.Context, sessionID string) (bool, int, time.Time, error) {
	key := rl.RateLimitKey(sessionID)

	// Try to get the current counter
	var counter RequestCounter
	err := database.GetCache(key, &counter)

	now := time.Now()
	var resetTime time.Time

	// If key doesn't exist or has expired, start a new window
	if err != nil {
		resetTime = now.Add(rl.windowDuration)
		counter = RequestCounter{
			Count:     1,
			ResetTime: resetTime,
		}
		if storeErr := database.SetCache(key, counter, rl.windowDuration); storeErr != nil {
			return false, 0, time.Time{}, fmt.Errorf("failed to set rate limit counter: %w", storeErr)
		}

		remaining := rl.requestsPerWindow - 1
		return true, remaining, resetTime, nil
	}

	// Check if window has expired
	if now.After(counter.ResetTime) {
		resetTime = now.Add(rl.windowDuration)
		counter = RequestCounter{
			Count:     1,
			ResetTime: resetTime,
		}
		if storeErr := database.SetCache(key, counter, rl.windowDuration); storeErr != nil {
			return false, 0, time.Time{}, fmt.Errorf("failed to reset rate limit counter: %w", storeErr)
		}

		remaining := rl.requestsPerWindow - 1
		return true, remaining, resetTime, nil
	}

	// Window still active - check if limit exceeded
	if counter.Count >= rl.requestsPerWindow {
		remaining := 0
		return false, remaining, counter.ResetTime, nil
	}

	// Increment counter and update in Redis
	counter.Count++
	timeToExpire := counter.ResetTime.Sub(now)
	if timeToExpire < 0 {
		timeToExpire = rl.windowDuration
	}

	if storeErr := database.SetCache(key, counter, timeToExpire); storeErr != nil {
		return false, 0, time.Time{}, fmt.Errorf("failed to increment rate limit counter: %w", storeErr)
	}

	remaining := rl.requestsPerWindow - counter.Count
	return true, remaining, counter.ResetTime, nil
}

// ResetSessionLimit manually resets the rate limit for a session (useful for admin operations)
func (rl *RateLimiterService) ResetSessionLimit(sessionID string) error {
	key := rl.RateLimitKey(sessionID)
	return database.DeleteCache(key)
}
