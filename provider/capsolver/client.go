// Package capsolver provides the CapSolver implementation of the unicap
// provider interface.
//
// CapSolver solves recognition tasks synchronously, so their results come
// back ready from CreateTask without polling.
package capsolver

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
)

// Name is the provider identifier.
const Name = "capsolver"

// API roots. StableBaseURL is CapSolver's secondary server.
const (
	BaseURL       = "https://api.capsolver.com"
	StableBaseURL = "https://api-stable.capsolver.com"
)

var (
	_ unicap.Provider       = (*Client)(nil)
	_ unicap.BalanceChecker = (*Client)(nil)
	_ unicap.Reporter       = (*Client)(nil)
	_ unicap.TaskSupporter  = (*Client)(nil)
	_ unicap.PollingAdvisor = (*Client)(nil)
)

// Client is a CapSolver provider. It is safe for concurrent use.
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

// WithBaseURL overrides the API root, for example with StableBaseURL.
func WithBaseURL(baseURL string) Option {
	return solverapi.WithBaseURL(baseURL)
}

// WithAppID sets the developer appId from the CapSolver dashboard, which
// credits usage to your application.
func WithAppID(appID string) Option {
	return solverapi.WithAppID(appID)
}

// WithCallbackURL sets the URL CapSolver posts results to when a task
// finishes.
func WithCallbackURL(callbackURL string) Option {
	return solverapi.WithCallbackURL(callbackURL)
}

// New creates a CapSolver provider.
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

// polling follows CapSolver's guidance: poll about every 3 seconds, at most
// 120 times, within the 5 minutes results are kept.
var polling = unicap.Polling{
	InitialDelay:       time.Second,
	Interval:           3 * time.Second,
	MaxInterval:        3 * time.Second,
	Multiplier:         1,
	Timeout:            5 * time.Minute,
	MaxTransientErrors: 3,
}

// errorRules classifies CapSolver's documented error codes.
var errorRules = solverapi.ErrorRules{
	"ERROR_KEY_DENIED_ACCESS":   unicap.ErrInvalidAPIKey,
	"ERROR_ZERO_BALANCE":        unicap.ErrInsufficientFunds,
	"ERROR_SETTLEMENT_FAILED":   unicap.ErrInsufficientFunds,
	"ERROR_TASKID_INVALID":      unicap.ErrTaskNotFound,
	"ERROR_INVALID_TASK_DATA":   unicap.ErrInvalidTask,
	"ERROR_UNKNOWN_QUESTION":    unicap.ErrInvalidTask,
	"ERROR_INVALID_IMAGE":       unicap.ErrInvalidTask,
	"ERROR_PARSE_IMAGE_FAIL":    unicap.ErrInvalidTask,
	"ERROR_TASK_NOT_SUPPORTED":  unicap.ErrUnsupportedTask,
	"ERROR_CAPTCHA_UNSOLVABLE":  unicap.ErrUnsolvable,
	"ERROR_TASK_TIMEOUT":        unicap.ErrUnsolvable,
	"ERROR_SERVICE_UNAVALIABLE": unicap.ErrUnavailable, //nolint:misspell // CapSolver spells the code this way.
	"ERROR_RATE_LIMIT":          unicap.ErrUnavailable,
	"ERROR_PROXY_BANNED":        unicap.ErrProxy,
	"ERROR_IP_BANNED":           unicap.ErrIPBlocked,
	"ERROR_KEY_TEMP_BLOCKED":    unicap.ErrIPBlocked,
}

// mappers registers every task type CapSolver solves.
var mappers = solverapi.Mappers{
	unicap.TaskTypeReCaptchaV2:             solverapi.Map(mapReCaptchaV2),
	unicap.TaskTypeReCaptchaV2Enterprise:   solverapi.Map(mapReCaptchaV2Enterprise),
	unicap.TaskTypeReCaptchaV3:             solverapi.Map(mapReCaptchaV3),
	unicap.TaskTypeReCaptchaV3Enterprise:   solverapi.Map(mapReCaptchaV3Enterprise),
	unicap.TaskTypeReCaptchaClassification: solverapi.Map(mapReCaptchaClassification),
	unicap.TaskTypeTurnstile:               solverapi.Map(mapTurnstile),
	unicap.TaskTypeCloudflareChallenge:     solverapi.Map(mapCloudflareChallenge),
	unicap.TaskTypeDataDome:                solverapi.Map(mapDataDome),
	unicap.TaskTypeAWSWAF:                  solverapi.Map(mapAWSWAF),
	unicap.TaskTypeAWSWAFClassification:    solverapi.Map(mapAWSWAFClassification),
	unicap.TaskTypeBotDeflector:            solverapi.Map(mapBotDeflector),
	unicap.TaskTypeGeeTest:                 solverapi.Map(mapGeeTest),
	unicap.TaskTypeGeeTestV4:               solverapi.Map(mapGeeTestV4),
	unicap.TaskTypeMTCaptcha:               solverapi.Map(mapMTCaptcha),
	unicap.TaskTypeImageToText:             solverapi.Map(mapImageToText),
	unicap.TaskTypeVisionEngine:            solverapi.Map(mapVisionEngine),
}

// variant returns CapSolver's proxyless and proxy wire types for base.
func variant(base string) solverapi.Variant {
	return solverapi.Variant{
		Proxyless: base + "ProxyLess",
		Proxy:     base,
	}
}
