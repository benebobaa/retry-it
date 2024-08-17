package retryit

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

// Do executes the provided operation with retries based on the specified configuration.
// The operation is retried according to the max attempts, delay, and backoff strategy set in the config.
// If the operation succeeds, it returns nil. If all retry attempts fail, it returns the last error encountered.
func Do(ctx context.Context, op Operation, opts ...Option) error {
	// Start with the default configuration and apply any provided options.
	cfg := defaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	var err error
	delay := cfg.InitialDelay

	// Loop through the number of allowed attempts.
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		// Execute the operation and return if it succeeds.
		if err = op(ctx); err == nil {
			return nil
		}

		// If this was the last attempt, break out of the loop to return the error.
		if attempt == cfg.MaxAttempts-1 {
			break
		}

		// Calculate the delay before the next retry attempt.
		delay = calculateNextDelay(delay, cfg)

		// Wait for the delay period or exit if the context is canceled.
		select {
		case <-ctx.Done():
			// If the context is done, return the context's error.
			return ctx.Err()
		case <-time.After(delay):
			// Wait for the calculated delay before the next retry.
		}
	}

	// If all attempts fail, return an error indicating the maximum retries were reached.
	return errors.New("max retry attempts reached: " + err.Error())
}

// calculateNextDelay computes the delay before the next retry attempt based on the current delay,
// backoff factor, maximum delay, and optional jitter.
// It increases the delay exponentially according to the backoff factor and applies jitter if configured.
func calculateNextDelay(currentDelay time.Duration, cfg Config) time.Duration {
	// Calculate the next delay using the backoff factor.
	nextDelay := time.Duration(float64(currentDelay) * cfg.BackoffFactor)

	// Ensure the delay does not exceed the maximum allowed delay.
	if nextDelay > cfg.MaxDelay {
		nextDelay = cfg.MaxDelay
	}

	// Apply jitter to add randomness to the delay.
	if cfg.Jitter > 0 {
		jitter := rand.Float64() * cfg.Jitter
		nextDelay = time.Duration(float64(nextDelay) * (1 + jitter))
	}

	return nextDelay
}
