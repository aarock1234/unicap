// Package twocaptcha provides the 2Captcha implementation of the unicap
// provider interface.
package twocaptcha

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aarock1234/unicap"
	"github.com/aarock1234/unicap/internal/solverapi"
)

const (
	name    = "2captcha"
	baseURL = "https://api.2captcha.com"
)

// Option configures the provider.
type Option = solverapi.Option

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(h *http.Client) Option {
	return solverapi.WithHTTPClient(h)
}

// WithBaseURL sets a custom base URL. Intended for testing.
func WithBaseURL(u string) Option {
	return solverapi.WithBaseURL(u)
}

// WithLogger sets a custom logger.
func WithLogger(l *slog.Logger) Option {
	return solverapi.WithLogger(l)
}

// New creates a 2Captcha provider.
func New(apiKey string, opts ...Option) (unicap.Provider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key: %w", unicap.ErrInvalidAPIKey)
	}

	errs := solverapi.StandardErrorMapper(
		name,
		[]string{"ERROR_KEY_DOES_NOT_EXIST", "ERROR_WRONG_USER_KEY"},
		[]string{"ERROR_ZERO_BALANCE"},
		[]string{"ERROR_TASK_ABSENT"},
		[]string{"ERROR_WRONG_TASK_DATA"},
	)

	return solverapi.New(name, baseURL, apiKey, mapTask, errs, opts...), nil
}
