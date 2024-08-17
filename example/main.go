package main

import (
	"context"
	"fmt"

	retryit "github.com/benebobaa/retry-it"
)

func main() {
	ctx := context.Background()

	counter := 0
	operation := func(ctx context.Context) error {
		// Your operation logic here
		// Return nil if successful, or an error if it needs to be retried

		counter++
		fmt.Println("Operation attempt:", counter)

		return fmt.Errorf("operation failed")
	}

	err := retryit.Do(ctx, operation)

	if err != nil {
		fmt.Println("All retries failed:", err)
	}
}
