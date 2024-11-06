package util

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/cenkalti/backoff/v4"
)

const (
	initialBackoff = 200 * time.Millisecond // Initial backoff duration
	maxBackoff     = 30 * time.Second       // Maximum backoff duration
)

var (
	ErrFailedAfterNAttempts = errors.New("operation failed after N attempts")
)

// RetryWithBackoff retries the operation with exponential backoff and a maximum number of attempts.
func RetryWithBackoff(ctx context.Context, maxAttempts uint64, log *slog.Logger, operation func() error) error {
	// Create and configure the ExponentialBackOff instance
	eb := backoff.NewExponentialBackOff()
	eb.InitialInterval = initialBackoff
	eb.MaxInterval = maxBackoff

	b := backoff.WithMaxRetries(eb, maxAttempts)

	err := backoff.Retry(func() error {
		select {
		case <-ctx.Done():
			log.Info("Context canceled, stopping retries.")
			return backoff.Permanent(ctx.Err())
		default:
			err := operation()
			if err != nil {
				// Log and retry unless it's a permanent error
				log.Warn("Operation failed, will retry", "error", err)
				return err
			}
			return nil // Success
		}
	}, backoff.WithContext(b, ctx))

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return err
		}
		log.Error("Operation failed after max attempts", "error", err)
		return ErrFailedAfterNAttempts
	}
	return nil
}

func RetryWithInfiniteBackoff(ctx context.Context, log *slog.Logger, operation func() error) error {
    eb := backoff.NewExponentialBackOff()
    eb.InitialInterval = initialBackoff
    eb.MaxInterval = maxBackoff
    eb.MaxElapsedTime = 0 // Infinite retry

    err := backoff.Retry(func() error {
        select {
        case <-ctx.Done():
            log.Info("Context canceled, stopping retries.")
            return backoff.Permanent(ctx.Err())
        default:
            err := operation()
            if err != nil {
                // Log and retry unless it's a permanent error
                log.Warn("Operation failed, will retry", "error", err)
                return err
            }
            return nil // Success
        }
    }, backoff.WithContext(eb, ctx))

    if err != nil {
        return fmt.Errorf("operation failed: %w", err)
    }
    return nil
}
