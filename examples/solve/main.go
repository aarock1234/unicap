// Package main solves Google's reCAPTCHA v2 demo with a provider chosen at
// runtime, then checks the balance and reports the result.
//
// Set CAPSOLVER_API_KEY, TWOCAPTCHA_API_KEY, or ANTICAPTCHA_API_KEY, then run:
//
//	go run ./examples/solve -provider capsolver
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/provider"
	"github.com/aarock1234/unicap/v3/tasks"
)

// apiKeyVariables maps provider names to the environment variable holding
// their API key.
var apiKeyVariables = map[string]string{
	"capsolver":   "CAPSOLVER_API_KEY",
	"2captcha":    "TWOCAPTCHA_API_KEY",
	"anticaptcha": "ANTICAPTCHA_API_KEY",
}

func main() {
	if err := run(); err != nil {
		slog.Error("solve example failed", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	providerName := flag.String("provider", "capsolver", "provider name: capsolver, 2captcha, or anticaptcha")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	client, err := newClient(*providerName, logger)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	balance, err := client.Balance(ctx)
	if err != nil {
		return fmt.Errorf("checking balance: %w", err)
	}

	logger.InfoContext(ctx, "account balance",
		slog.Float64("usd", balance),
	)

	result, err := client.Solve(ctx, &tasks.ReCaptchaV2Task{
		WebsiteURL: "https://www.google.com/recaptcha/api2/demo",
		WebsiteKey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
	})

	switch {
	case errors.Is(err, unicap.ErrInsufficientFunds):
		return errors.New("top up the provider account and try again")
	case err != nil:
		return fmt.Errorf("solving captcha: %w", err)
	}

	logger.InfoContext(ctx, "captcha solved",
		slog.String("task_id", result.TaskID),
		slog.Float64("cost_usd", result.Cost),
		slog.Duration("elapsed", result.Elapsed()),
		slog.Int("token_length", len(result.Solution.Token)),
	)

	// After submitting the token to the target site, tell the provider
	// whether it worked.
	if err := client.Report(ctx, result, true); err != nil && !errors.Is(err, unicap.ErrUnsupportedOperation) {
		return fmt.Errorf("reporting result: %w", err)
	}

	return nil
}

// newClient creates a unicap client for the named provider using its API key
// from the environment.
func newClient(providerName string, logger *slog.Logger) (*unicap.Client, error) {
	variable, ok := apiKeyVariables[providerName]
	if !ok {
		return nil, fmt.Errorf("%q: %w", providerName, provider.ErrUnknownProvider)
	}

	apiKey := os.Getenv(variable)
	if apiKey == "" {
		return nil, fmt.Errorf("%s is not set", variable)
	}

	solver, err := provider.NewRegistry().New(providerName, apiKey,
		provider.WithLogger(logger),
	)
	if err != nil {
		return nil, fmt.Errorf("creating provider: %w", err)
	}

	client, err := unicap.New(solver,
		unicap.WithLogger(logger),
		unicap.WithSolveRetries(1),
	)
	if err != nil {
		return nil, fmt.Errorf("creating client: %w", err)
	}

	return client, nil
}
