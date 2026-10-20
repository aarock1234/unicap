// Package main demonstrates plugging a custom provider into unicap.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/tasks"
)

func main() {
	if err := run(); err != nil {
		slog.Error("custom provider example failed", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	provider, err := newCustomProvider("your-api-key")
	if err != nil {
		return fmt.Errorf("creating provider: %w", err)
	}

	client, err := unicap.New(provider)
	if err != nil {
		return fmt.Errorf("creating client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	result, err := client.Solve(ctx, &tasks.ReCaptchaV2Task{
		WebsiteURL: "https://example.com",
		WebsiteKey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
	})
	if err != nil {
		return fmt.Errorf("solving captcha: %w", err)
	}

	fmt.Printf("token: %s\n", result.Solution.Token)

	return nil
}
