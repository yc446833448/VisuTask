package agent

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
)

// RetryPolicy implements exponential backoff with configurable parameters
type RetryPolicy struct {
	InitialDelay  time.Duration
	MaxDelay      time.Duration
	BackoffFactor float64
	MaxRetries    int
}

func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		InitialDelay:  2 * time.Second,
		MaxDelay:      30 * time.Second,
		BackoffFactor: 2.0,
		MaxRetries:    5,
	}
}

// ShouldRetry determines if an error is retryable
func (r *RetryPolicy) ShouldRetry(err error) bool {
	if err == nil {
		return false
	}

	msg := err.Error()

	// Not retryable: context overflow, auth failures
	if strings.Contains(msg, "context_overflow") ||
		strings.Contains(msg, "authentication") ||
		strings.Contains(msg, "invalid_api_key") ||
		strings.Contains(msg, "401") ||
		strings.Contains(msg, "403") {
		return false
	}

	// Retryable: server errors, rate limits, network issues
	if strings.Contains(msg, "500") ||
		strings.Contains(msg, "502") ||
		strings.Contains(msg, "503") ||
		strings.Contains(msg, "529") ||
		strings.Contains(msg, "429") ||
		strings.Contains(msg, "too_many_requests") ||
		strings.Contains(msg, "rate_limit") ||
		strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "connection") {
		return true
	}

	return false
}

// Wait performs exponential backoff wait, respecting context cancellation
func (r *RetryPolicy) Wait(ctx context.Context, attempt int, err error) error {
	delay := r.CalculateDelay(attempt, err)

	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// CalculateDelay computes the wait time for a given attempt
func (r *RetryPolicy) CalculateDelay(attempt int, err error) time.Duration {
	// Check for retry-after header in the error
	if retryAfter := extractRetryAfter(err); retryAfter > 0 {
		return retryAfter
	}

	// Exponential backoff
	delay := float64(r.InitialDelay) * math.Pow(r.BackoffFactor, float64(attempt))
	if delay > float64(r.MaxDelay) {
		delay = float64(r.MaxDelay)
	}
	return time.Duration(delay)
}

// extractRetryAfter attempts to parse retry-after from error
func extractRetryAfter(err error) time.Duration {
	// In a real implementation, this would check HTTP response headers
	// For now, return 0 to use exponential backoff
	return 0
}

// IsRetryAfterHeader checks if the error contains a retry-after HTTP header
func IsRetryAfterHeader(resp *http.Response) (time.Duration, bool) {
	if resp == nil {
		return 0, false
	}

	// Check retry-after header (seconds)
	if ra := resp.Header.Get("Retry-After"); ra != "" {
		// Try parsing as seconds
		var seconds int
		if _, err := fmt.Sscanf(ra, "%d", &seconds); err == nil {
			return time.Duration(seconds) * time.Second, true
		}
		// Try parsing as HTTP date
		if t, err := http.ParseTime(ra); err == nil {
			d := time.Until(t)
			if d > 0 {
				return d, true
			}
		}
	}

	// Check retry-after-ms header
	if raMs := resp.Header.Get("Retry-After-Ms"); raMs != "" {
		var ms int
		if _, err := fmt.Sscanf(raMs, "%d", &ms); err == nil {
			return time.Duration(ms) * time.Millisecond, true
		}
	}

	return 0, false
}
