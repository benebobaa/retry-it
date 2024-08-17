package retryit

import (
	"time"
)

// WithMaxAttempts configures the maximum number of retry attempts.
// If the operation fails, it will be retried up to the specified number of times.
// A value of 0 or less means there will be no retries.
func WithMaxAttempts(attempts int) Option {
	return func(c *Config) {
		c.MaxAttempts = attempts
	}
}

// WithInitialDelay sets the initial delay before the first retry attempt.
// The delay is specified as a time.Duration value and is applied after the first failure.
// This delay will be adjusted for subsequent retries based on the backoff factor.
func WithInitialDelay(delay time.Duration) Option {
	return func(c *Config) {
		c.InitialDelay = delay
	}
}

// WithMaxDelay configures the maximum delay between retry attempts.
// No matter how many retries occur or the backoff strategy used,
// the delay between retries will not exceed this maximum duration.
func WithMaxDelay(delay time.Duration) Option {
	return func(c *Config) {
		c.MaxDelay = delay
	}
}

// WithBackoffFactor sets the backoff factor used to increase the delay between retries.
// A factor greater than 1 will cause the delay to grow exponentially, helping to
// space out retry attempts further as failures continue. A factor of 1 results
// in a constant delay between attempts.
func WithBackoffFactor(factor float64) Option {
	return func(c *Config) {
		c.BackoffFactor = factor
	}
}

// WithJitter adds a randomness factor to the retry delays to avoid synchronized retries.
// The jitter is applied as a percentage of the calculated delay to randomize it slightly.
// This helps to prevent the "thundering herd" problem in distributed systems.
func WithJitter(jitter float64) Option {
	return func(c *Config) {
		c.Jitter = jitter
	}
}
