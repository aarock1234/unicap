// Package main demonstrates solving a captcha with the built-in CapSolver
// provider.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/aarock1234/unicap"
	"github.com/aarock1234/unicap/provider/capsolver"
	"github.com/aarock1234/unicap/tasks"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal error", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("loading .env: %w", err)
	}

	provider, err := capsolver.New(os.Getenv("CAPSOLVER_API_KEY"))
	if err != nil {
		return fmt.Errorf("creating provider: %w", err)
	}

	client, err := unicap.New(provider)
	if err != nil {
		return fmt.Errorf("creating client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Solve automatically polls until the solution is ready.
	solution, err := client.Solve(ctx, &tasks.ReCaptchaV2Task{
		WebsiteURL: "https://www.google.com/recaptcha/api2/demo",
		WebsiteKey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
	})
	if err != nil {
		return fmt.Errorf("solving captcha: %w", err)
	}

	fmt.Printf("token: %s\n", solution.Token)

	return nil
}
