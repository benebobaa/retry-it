package retryit

import (
	"context"
	"time"
)

// Operation defines a function signature that represents an operation
// that can be retried. The function should return an error to indicate
// a failure, triggering a retry based on the retry configuration.
type Operation func(context.Context) error

// Config encapsulates the configuration options for the retry mechanism.
// This structure allows fine-tuning of how retries are handled, including
// the maximum number of attempts, delays between retries, and backoff strategy.
type Config struct {
	// MaxAttempts specifies the maximum number of attempts to retry the operation.
	// A value of 0 or less means there will be no retries.
	MaxAttempts int

	// InitialDelay is the duration to wait before the first retry attempt.
	// This delay will be increased based on the backoff factor for each subsequent retry.
	InitialDelay time.Duration

	// MaxDelay defines the maximum delay duration between retries.
	// The delay will never exceed this value, even with a backoff strategy applied.
	MaxDelay time.Duration

	// BackoffFactor determines the multiplier applied to the delay between retries.
	// A factor greater than 1 will result in exponential backoff, while a value of 1
	// will maintain a constant delay between attempts.
	BackoffFactor float64

	// Jitter adds randomness to the delay between retries to avoid
	// the "thundering herd" problem when multiple operations are retried simultaneously.
	// The jitter value is applied as a percentage of the current delay.
	Jitter float64
}

// Option defines a function type used to modify the retry configuration.
// It allows for functional options to be applied when initializing the Config struct,
// enabling customization of retry behavior.
type Option func(*Config)

// defaultConfig provides a set of well-chosen default values for the retry configuration.
// These defaults aim to provide a balanced retry mechanism, suitable for a wide range of use cases.
// The default configuration includes a maximum of 5 attempts, an initial delay of 100ms,
// a maximum delay of 30 seconds, an exponential backoff factor of 2, and a 10% jitter.
var defaultConfig = Config{
	MaxAttempts:   5,
	InitialDelay:  100 * time.Millisecond,
	MaxDelay:      30 * time.Second,
	BackoffFactor: 2,
	Jitter:        0.1,
}
