// Package noptcha provides the Noptcha implementation of the unicap provider
// interface, using Noptcha's native v1 API.
//
// Noptcha describes itself as pre-launch: its prices are unpublished, its
// terms are a draft, and its documented limits are starting values.
//
// Noptcha solves some task types synchronously, so their results come back
// ready from CreateTask without polling.
package noptcha

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/httpapi"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
)

// Name is the provider identifier.
const Name = "noptcha"

// BaseURL is the API root.
const BaseURL = "https://api.noptcha.com"

var (
	_ unicap.Provider       = (*Client)(nil)
	_ unicap.BalanceChecker = (*Client)(nil)
	_ unicap.Reporter       = (*Client)(nil)
	_ unicap.TaskSupporter  = (*Client)(nil)
	_ unicap.PollingAdvisor = (*Client)(nil)
)

// Client is a Noptcha provider. It is safe for concurrent use.
type Client struct {
	api *httpapi.Client
}

// Option configures a Client.
type Option = httpapi.Option

// WithHTTPClient sets the HTTP client used for API requests.
func WithHTTPClient(httpClient *http.Client) Option {
	return httpapi.WithHTTPClient(httpClient)
}

// WithLogger sets the logger for request diagnostics.
func WithLogger(logger *slog.Logger) Option {
	return httpapi.WithLogger(logger)
}

// WithBaseURL overrides the API root.
func WithBaseURL(baseURL string) Option {
	return httpapi.WithBaseURL(baseURL)
}

// New creates a Noptcha provider.
func New(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("%s api key is empty: %w", Name, unicap.ErrInvalidAPIKey)
	}

	config := httpapi.Config{
		Name:    Name,
		BaseURL: BaseURL,
		Header: http.Header{
			"Authorization": []string{
				"Bearer " + apiKey,
			},
		},
	}

	return &Client{
		api: httpapi.New(config, opts...),
	}, nil
}

// polling is unicap's choice: Noptcha prescribes no interval for the v1 API.
// One task may be polled 5 times a second, so a fixed 2 second pace stays far
// below the limit, and the 3 minute timeout outlasts the longest task
// deadline of 120 seconds.
var polling = unicap.Polling{
	InitialDelay:       2 * time.Second,
	Interval:           2 * time.Second,
	MaxInterval:        2 * time.Second,
	Multiplier:         1,
	Timeout:            3 * time.Minute,
	MaxTransientErrors: 3,
}

// requestErrorRules classifies the codes of Noptcha's request errors, which
// arrive as problem details with an HTTP error status. A nil sentinel marks
// a documented code that matches no unicap condition. The forbidden code
// covers a missing scope, a disallowed IP or type, and a suspended
// organisation alike, so it maps to the key being refused.
var requestErrorRules = map[string]error{
	"invalid_request":      unicap.ErrInvalidTask,
	"unsupported_type":     unicap.ErrUnsupportedTask,
	"unauthorized":         unicap.ErrInvalidAPIKey,
	"insufficient_balance": unicap.ErrInsufficientFunds,
	"budget_exceeded":      unicap.ErrInsufficientFunds,
	"forbidden":            unicap.ErrInvalidAPIKey,
	"not_found":            unicap.ErrTaskNotFound,
	"idempotency_conflict": nil,
	"task_not_solved":      nil,
	"request_too_large":    unicap.ErrInvalidTask,
	"invalid_params":       unicap.ErrInvalidTask,
	"rate_limited":         unicap.ErrUnavailable,
	"concurrency_limited":  unicap.ErrUnavailable,
	"internal":             unicap.ErrUnavailable,
	"service_unavailable":  unicap.ErrUnavailable,
}

// taskFailureRules classifies the error codes of failed tasks. A failed task
// is final, so no code maps to unicap.ErrUnavailable, which would make the
// poller ask for the same dead task again.
var taskFailureRules = map[string]error{
	"unsolvable":    unicap.ErrUnsolvable,
	"timeout":       unicap.ErrUnsolvable,
	"proxy_error":   unicap.ErrProxy,
	"invalid_site":  unicap.ErrInvalidTask,
	"invalid_image": unicap.ErrInvalidTask,
	"cancelled":     unicap.ErrTaskFailed,
}

// mappers registers every task type Noptcha solves, each under its native
// type. FunCaptcha is left out because Noptcha does not support it at the
// moment.
var mappers = solverapi.Mappers{
	unicap.TaskTypeReCaptchaV2:             solverapi.Map(native("recaptcha.v2", mapReCaptchaV2)),
	unicap.TaskTypeReCaptchaV2Enterprise:   solverapi.Map(native("recaptcha.v2.enterprise", mapReCaptchaV2Enterprise)),
	unicap.TaskTypeReCaptchaV3:             solverapi.Map(native("recaptcha.v3", mapReCaptchaV3)),
	unicap.TaskTypeReCaptchaV3Enterprise:   solverapi.Map(native("recaptcha.v3.enterprise", mapReCaptchaV3Enterprise)),
	unicap.TaskTypeReCaptchaClassification: solverapi.Map(native("image.grid.recaptcha", mapReCaptchaClassification)),
	unicap.TaskTypeHCaptcha:                solverapi.Map(native("hcaptcha", mapHCaptcha)),
	unicap.TaskTypeTurnstile:               solverapi.Map(native("turnstile", mapTurnstile)),
	unicap.TaskTypeCloudflareChallenge:     solverapi.Map(native("cloudflare.challenge", mapCloudflareChallenge)),
	unicap.TaskTypeImageToText:             solverapi.Map(native("image.text", mapImageToText)),
}

// native adapts a function that maps a task to its params object into one
// that returns the whole create-task body for the native task type.
func native[T unicap.Task, P any](taskType string, mapParams func(T) (P, error)) func(T) (taskRequest[P], error) {
	return func(task T) (taskRequest[P], error) {
		params, err := mapParams(task)
		if err != nil {
			return taskRequest[P]{}, err
		}

		return taskRequest[P]{
			Type:   taskType,
			Params: params,
		}, nil
	}
}
