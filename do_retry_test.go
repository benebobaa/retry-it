package retryit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDo(t *testing.T) {
	t.Run("Successful operation", func(t *testing.T) {
		attempts := 0
		err := Do(context.Background(), func(ctx context.Context) error {
			attempts++
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, attempts)
	})

	t.Run("Retry until success", func(t *testing.T) {
		attempts := 0
		err := Do(context.Background(), func(ctx context.Context) error {
			attempts++
			if attempts < 3 {
				return errors.New("temporary error")
			}
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 3, attempts)
	})

	t.Run("Max attempts reached", func(t *testing.T) {
		attempts := 0
		err := Do(context.Background(), func(ctx context.Context) error {
			attempts++
			return errors.New("persistent error")
		})
		assert.Error(t, err)
		assert.Equal(t, defaultConfig.MaxAttempts, attempts)
	})

	t.Run("Context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		attempts := 0
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()
		err := Do(ctx, func(ctx context.Context) error {
			attempts++
			time.Sleep(100 * time.Millisecond)
			return errors.New("slow operation")
		})
		assert.Error(t, err)
		assert.True(t, errors.Is(err, context.Canceled))
		assert.Less(t, attempts, defaultConfig.MaxAttempts)
	})
}

func TestOptions(t *testing.T) {
	t.Run("WithMaxAttempts", func(t *testing.T) {
		attempts := 0
		err := Do(context.Background(),
			func(ctx context.Context) error {
				attempts++
				return errors.New("error")
			},
			WithMaxAttempts(3),
		)
		assert.Error(t, err)
		assert.Equal(t, 3, attempts)
	})

	t.Run("WithInitialDelay", func(t *testing.T) {
		start := time.Now()
		_ = Do(context.Background(),
			func(ctx context.Context) error {
				return errors.New("error")
			},
			WithMaxAttempts(2),
			WithInitialDelay(100*time.Millisecond),
		)
		duration := time.Since(start)
		assert.GreaterOrEqual(t, duration, 100*time.Millisecond)
	})

	t.Run("WithMaxDelay", func(t *testing.T) {
		start := time.Now()
		_ = Do(context.Background(),
			func(ctx context.Context) error {
				return errors.New("error")
			},
			WithMaxAttempts(3),
			WithInitialDelay(100*time.Millisecond),
			WithMaxDelay(150*time.Millisecond),
		)
		duration := time.Since(start)
		assert.Less(t, duration, 500*time.Millisecond)
	})

	t.Run("WithBackoffFactor", func(t *testing.T) {
		delays := []time.Duration{}
		start := time.Now()
		_ = Do(context.Background(),
			func(ctx context.Context) error {
				delays = append(delays, time.Since(start))
				start = time.Now()
				return errors.New("error")
			},
			WithMaxAttempts(3),
			WithInitialDelay(100*time.Millisecond),
			WithBackoffFactor(2),
		)
		require.Len(t, delays, 3)
		t.Logf("Delays: %v", delays)

		// Check that the ratio between delays is approximately 2
		ratio := float64(delays[2]) / float64(delays[1])
		assert.InDelta(t, 2.0, ratio, 0.5, "The ratio between delays should be approximately 2")
	})

	t.Run("WithJitter", func(t *testing.T) {
		delays := []time.Duration{}
		start := time.Now()
		_ = Do(context.Background(),
			func(ctx context.Context) error {
				delay := time.Since(start)
				delays = append(delays, delay)
				t.Logf("Delay: %v", delay)
				start = time.Now()
				return errors.New("error")
			},
			WithMaxAttempts(3),
			WithInitialDelay(100*time.Millisecond),
			WithJitter(0.5),
		)
		require.Len(t, delays, 3)
		for i, delay := range delays[1:] {
			assert.GreaterOrEqual(t, delay.Milliseconds(), int64(50), "Delay %d too short: %v", i+1, delay)
			assert.Less(t, delay.Milliseconds(), int64(600), "Delay %d too long: %v", i+1, delay)
		}
	})
}
