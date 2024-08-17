package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	retryit "github.com/benebobaa/retry-it"
)

func main() {
	ctx := context.Background()

	err := retryit.Do(ctx,
		func(ctx context.Context) error {
			return makeRequest(ctx, "https://api.example.com/data")
		},
		retryit.WithMaxAttempts(5),
		retryit.WithInitialDelay(100*time.Millisecond),
		retryit.WithMaxDelay(2*time.Second),
		retryit.WithBackoffFactor(2),
		retryit.WithJitter(0.1),
	)

	if err != nil {
		fmt.Println("Failed to make request after retries:", err)
	} else {
		fmt.Println("Request successful")
	}
}

func makeRequest(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 || resp.StatusCode >= 500 {
		return fmt.Errorf("received status code %d", resp.StatusCode)
	}

	// Process successful response here
	fmt.Printf("Request successful with status code: %d\n", resp.StatusCode)
	return nil
}
