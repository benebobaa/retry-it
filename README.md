
# Retry-It

`retry-it` is a Go library that provides a simple and flexible retry mechanism for handling transient failures in your applications.

## Installation

To install `retry-it`, use `go get`:

```bash
go get github.com/benebobaa/retry-it
```

## Usage

Here's a basic example of how to use `retry-it`:

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/benebobaa/retry-it"
)

func main() {
    ctx := context.Background()

    err := retryit.Do(ctx, func(ctx context.Context) error {
        // Your operation here
        return fmt.Errorf("temporary error")
    })

    if err != nil {
        fmt.Println("Operation failed after retries:", err)
    }
}
```

## Retry

### func Do

```go
func Do(ctx context.Context, op Operation, opts ...Option) error
```

`Do` executes the given operation with retries based on the provided configuration. It returns an error if the operation fails after all retry attempts.

#### Parameters:

- `ctx context.Context`: The context for the operation.
- `op Operation`: The operation to be retried. It should be a function that returns an error.
- `opts ...Option`: Optional configuration options.

### type Operation

```go
type Operation func(context.Context) error
```

`Operation` represents a function that can be retried. It takes a context and returns an error.

### Options

`retry-it` provides several options to configure the retry behavior:

```go
WithMaxAttempts(attempts int) Option
WithInitialDelay(delay time.Duration) Option
WithMaxDelay(delay time.Duration) Option
WithBackoffFactor(factor float64) Option
WithJitter(jitter float64) Option
```

#### Example with options:

```go
err := retry.Do(ctx, operation,
    retry.WithMaxAttempts(5),
    retry.WithInitialDelay(100*time.Millisecond),
    retry.WithMaxDelay(2*time.Second),
    retry.WithBackoffFactor(2.0),
    retry.WithJitter(0.1),
)
```

## Default Configuration

If no options are provided, `retry-it` uses the following default configuration:

- Max Attempts: 5
- Initial Delay: 100ms
- Max Delay: 30s
- Backoff Factor: 2.0
- Jitter: 0.1 (10%)

## Features

- Context support for cancellation and timeouts
- Exponential backoff with jitter
- Configurable retry attempts, delays, and backoff factors
- Easy to use with sensible defaults

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

If you have any questions or issues, please open an issue on the GitHub repository.
