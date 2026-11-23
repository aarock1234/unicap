// Package capshot provides the Capshot implementation of the unicap provider
// interface, using Capshot's native API.
//
// The provider covers only the task types Capshot documents in full. hCaptcha
// and Cloudflare are left out because the vendor's documentation does not
// specify them. For the same reason only AWS WAF tasks accept a proxy: it is
// the one family whose proxy task type the documentation names.
//
// Every Capshot task is asynchronous: CreateTask returns a pending task and
// the solution comes from GetTaskResult.
package capshot

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
)

// Name is the provider identifier.
const Name = "capshot"

// BaseURL is the API root.
const BaseURL = "https://api.capshot.net"

var (
	_ unicap.Provider       = (*Client)(nil)
	_ unicap.BalanceChecker = (*Client)(nil)
	_ unicap.TaskSupporter  = (*Client)(nil)
	_ unicap.PollingAdvisor = (*Client)(nil)
)

// Client is a Capshot provider. It is safe for concurrent use.
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

// New creates a Capshot provider.
func New(apiKey string, opts ...Option) (*Client, error) {
	config := solverapi.Config{
		Name:           Name,
		BaseURL:        BaseURL,
		APIKey:         apiKey,
		TaskIDFormat:   solverapi.TaskIDNumber,
		Mappers:        mappers,
		Errors:         errorRules,
		Polling:        polling,
		RefineSolution: refineSolution,
	}

	api, err := solverapi.New(config, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		api: api,
	}, nil
}

// polling follows Capshot's guidance: "Recommended polling interval: 1–2
// seconds." Capshot documents no time limit, so the 2 minute timeout is
// unicap's choice; typical solves are quoted at a few seconds.
var polling = unicap.Polling{
	InitialDelay:       time.Second,
	Interval:           2 * time.Second,
	MaxInterval:        2 * time.Second,
	Multiplier:         1,
	Timeout:            2 * time.Minute,
	MaxTransientErrors: 3,
}

// errorRules classifies Capshot's documented error codes. ERROR_TASK_PAYLOAD
// also covers an invalid proxy, which Capshot does not report separately.
var errorRules = solverapi.ErrorRules{
	"ERROR_KEY_DOES_NOT_EXIST": unicap.ErrInvalidAPIKey,
	"ERROR_NO_FUNDS":           unicap.ErrInsufficientFunds,
	"ERROR_CAPTCHA_UNSOLVABLE": unicap.ErrUnsolvable,
	"ERROR_NO_SUCH_METHOD":     unicap.ErrInvalidTask,
	"ERROR_TASK_NOT_SUPPORTED": unicap.ErrUnsupportedTask,
	"ERROR_TASK_NOT_FOUND":     unicap.ErrTaskNotFound,
	"ERROR_TASK_PAYLOAD":       unicap.ErrInvalidTask,
}

// mappers registers every task type Capshot documents in full.
var mappers = solverapi.Mappers{
	unicap.TaskTypeReCaptchaV2:           solverapi.Map(mapReCaptchaV2),
	unicap.TaskTypeReCaptchaV2Enterprise: solverapi.Map(mapReCaptchaV2Enterprise),
	unicap.TaskTypeReCaptchaV3:           solverapi.Map(mapReCaptchaV3),
	unicap.TaskTypeReCaptchaV3Enterprise: solverapi.Map(mapReCaptchaV3Enterprise),
	unicap.TaskTypeGeeTestV4:             solverapi.Map(mapGeeTestV4),
	unicap.TaskTypeFunCaptcha:            solverapi.Map(mapFunCaptcha),
	unicap.TaskTypeFriendlyCaptcha:       solverapi.Map(mapFriendlyCaptcha),
	unicap.TaskTypeAWSWAF:                solverapi.Map(mapAWSWAF),
}
