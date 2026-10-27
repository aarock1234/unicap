// Package solvecha provides the Solvecha implementation of the unicap
// provider interface.
//
// Solvecha recognises pictures and text only. It solves ImageToTextTask and
// ReCaptchaClassificationTask; token tasks such as reCAPTCHA, hCaptcha or
// Turnstile are not supported and fail with unicap.ErrUnsupportedTask.
// Results usually come back ready from CreateTask; a task that takes longer
// is fetched with GetTaskResult.
//
// See https://solvecha.com/docs for the API documentation.
package solvecha

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
)

// Name is the provider identifier.
const Name = "solvecha"

// BaseURL is the Solvecha API root.
const BaseURL = "https://api.solvecha.com"

var (
	_ unicap.Provider       = (*Client)(nil)
	_ unicap.BalanceChecker = (*Client)(nil)
	_ unicap.TaskSupporter  = (*Client)(nil)
	_ unicap.PollingAdvisor = (*Client)(nil)
)

// Client is a Solvecha provider. It is safe for concurrent use.
type Client struct {
	api *solverapi.Client
}

// Option configures a Client.
type Option = solverapi.Option

// WithHTTPClient sets the HTTP client used for API requests.
func WithHTTPClient(httpClient *http.Client) Option {
	return solverapi.WithHTTPClient(httpClient)
}

// WithLogger sets the logger for request diagnostics.
func WithLogger(logger *slog.Logger) Option {
	return solverapi.WithLogger(logger)
}

// WithBaseURL overrides the API root.
func WithBaseURL(baseURL string) Option {
	return solverapi.WithBaseURL(baseURL)
}

// New creates a Solvecha provider.
func New(apiKey string, opts ...Option) (*Client, error) {
	config := solverapi.Config{
		Name:         Name,
		BaseURL:      BaseURL,
		APIKey:       apiKey,
		TaskIDFormat: solverapi.TaskIDString,
		Mappers:      mappers,
		Errors:       errorRules,
		Polling:      polling,
	}

	api, err := solverapi.New(config, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		api: api,
	}, nil
}

// polling is the pace for the rare task that createTask does not finish
// itself, since createTask waits up to 15 seconds for the answer.
var polling = unicap.Polling{
	InitialDelay:       time.Second,
	Interval:           2 * time.Second,
	MaxInterval:        2 * time.Second,
	Multiplier:         1,
	Timeout:            2 * time.Minute,
	MaxTransientErrors: 3,
}

// errorRules classifies Solvecha's error codes.
var errorRules = solverapi.ErrorRules{
	"ERROR_KEY_DENIED_ACCESS":  unicap.ErrInvalidAPIKey,
	"ERROR_ZERO_BALANCE":       unicap.ErrInsufficientFunds,
	"ERROR_TASKID_INVALID":     unicap.ErrTaskNotFound,
	"ERROR_INVALID_TASK_DATA":  unicap.ErrInvalidTask,
	"ERROR_TASK_NOT_SUPPORTED": unicap.ErrUnsupportedTask,
	"ERROR_CAPTCHA_UNSOLVABLE": unicap.ErrUnsolvable,
	"ERROR_RATE_LIMIT":         unicap.ErrUnavailable,
}

// mappers registers the picture and text task types Solvecha solves.
var mappers = solverapi.Mappers{
	unicap.TaskTypeReCaptchaClassification: solverapi.Map(mapReCaptchaClassification),
	unicap.TaskTypeImageToText:             solverapi.Map(mapImageToText),
}
