package redisapp

import (
	"context"
	"errors"
	"log"
	"math"
	"time"
)

const (
	initialBackoff = 200 * time.Millisecond // Initial backoff duration
	maxBackoff     = 30 * time.Second       // Maximum backoff duration
)

var (
	ErrFailedAfterNAttempts = errors.New("operation failed after N attempts")
)

// Backoff function implementing an exponential backoff strategy
func backoff(attempt int) time.Duration {
	backoff := float64(initialBackoff) * math.Pow(2, float64(attempt))
	if backoff > float64(maxBackoff) {
		backoff = float64(maxBackoff)
	}
	return time.Duration(backoff)
}

func retryWithInfiniteBackoff(ctx context.Context, operation func() (bool, error)) (bool, error) {
	for attempt := 0; ; attempt++ {
		select {
		case <-ctx.Done():
			log.Println("Context canceled, stopping retries.")
			return false, ctx.Err()
		default:
			success, err := operation()
			if success {
				return true, nil
			}

			if err != nil {
				log.Printf("Operation failed (attempt %d): %v.", attempt+1, err)
				return false, err
			} 
			
			log.Printf("Operation not successful (attempt %d). Retrying...", attempt+1)
			
			time.Sleep(backoff(attempt))
		}
	}
}

func retryWithLimitedAttempts(ctx context.Context, operation func() (bool, error), maxAttempts int) (bool, error) {
	for attempt := 0; attempt < maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			log.Println("Context canceled, stopping retries.")
			return false, ctx.Err()
		default:
			success, err := operation()
			if success {
				return true, nil
			}

			if err != nil {
				log.Printf("Operation failed (attempt %d): %v.", attempt+1, err)
				return false, err
			}

			log.Printf("Operation not successful (attempt %d). Retrying...", attempt+1)

			time.Sleep(backoff(attempt))
		}
	}

	log.Println("Max retry attempts reached, stopping retries.")
	return false, ErrFailedAfterNAttempts
}
