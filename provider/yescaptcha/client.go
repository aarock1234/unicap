// Package yescaptcha provides the YesCaptcha implementation of the unicap
// provider interface, using YesCaptcha's native API.
//
// YesCaptcha solves reCAPTCHA image classification synchronously, so its
// result comes back ready from CreateTask without polling.
//
// The account balance is kept in points rather than US dollars, so the
// provider does not implement unicap.BalanceChecker.
package yescaptcha

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
)

// Name is the provider identifier.
const Name = "yescaptcha"

// API roots. ChinaBaseURL is YesCaptcha's server for mainland China. Tasks
// are not shared between the two: a task must be queried on the root that
// created it.
const (
	BaseURL      = "https://api.yescaptcha.com"
	ChinaBaseURL = "https://cn.yescaptcha.com"
)

var (
	_ unicap.Provider       = (*Client)(nil)
	_ unicap.TaskSupporter  = (*Client)(nil)
	_ unicap.PollingAdvisor = (*Client)(nil)
)

// Client is a YesCaptcha provider. It is safe for concurrent use.
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

// WithBaseURL overrides the API root, for example with ChinaBaseURL.
func WithBaseURL(baseURL string) Option {
	return solverapi.WithBaseURL(baseURL)
}

// New creates a YesCaptcha provider.
func New(apiKey string, opts ...Option) (*Client, error) {
	config := solverapi.Config{
		Name:           Name,
		BaseURL:        BaseURL,
		APIKey:         apiKey,
		TaskIDFormat:   solverapi.TaskIDString,
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

// polling follows YesCaptcha's guidance: wait 5 seconds after creating a
// task, poll every 5 seconds because the server answers a repeat query
// within 5 seconds with "processing" unchanged, and give up after about 3
// minutes; most tasks time out server-side after about 110 seconds.
var polling = unicap.Polling{
	InitialDelay:       5 * time.Second,
	Interval:           5 * time.Second,
	MaxInterval:        5 * time.Second,
	Multiplier:         1,
	Timeout:            3 * time.Minute,
	MaxTransientErrors: 3,
}

// errorRules classifies YesCaptcha's documented error codes.
// ERROR_PAGE_NOT_FOUND is left unclassified: it reports a wrong API root,
// not a task condition.
var errorRules = solverapi.ErrorRules{
	"ERROR_REQUIRED_FIELDS":           unicap.ErrInvalidTask,
	"ERROR_REQUEST_FORMATS":           unicap.ErrInvalidTask,
	"ERROR_TASK_ABSENT":               unicap.ErrInvalidTask,
	"ERROR_TASK_NOT_SUPPORTED":        unicap.ErrUnsupportedTask,
	"ERROR_BAD_REQUEST":               unicap.ErrInvalidTask,
	"ERROR_ILLEGAL_PARAMETERS":        unicap.ErrInvalidTask,
	"ERROR_KEY_DOES_NOT_EXIST":        unicap.ErrInvalidAPIKey,
	"ERROR_KEY_EXPIRED":               unicap.ErrInvalidAPIKey,
	"ERROR_KEY_IP_NOT_ALLOWED":        unicap.ErrIPBlocked,
	"ERROR_KEY_QUOTA_EXCEEDED":        unicap.ErrInsufficientFunds,
	"ERROR_ZERO_BALANCE":              unicap.ErrInsufficientFunds,
	"ERROR_THREADS_MAXIMUM":           unicap.ErrUnavailable,
	"ERROR_IP_BLOCKED_5MIN":           unicap.ErrIPBlocked,
	"ERROR_NO_SLOT_AVAILABLE":         unicap.ErrUnavailable,
	"ERROR_TASKID_INVALID":            unicap.ErrTaskNotFound,
	"ERROR_TASK_TIMEOUT":              unicap.ErrUnsolvable,
	"ERROR_CAPTCHA_UNSOLVABLE":        unicap.ErrUnsolvable,
	"ERROR_SETTLEMENT_FAILED":         unicap.ErrInsufficientFunds,
	"ERROR_NO_SUCH_CAPCHA_ID":         unicap.ErrTaskNotFound,
	"ERROR_RECAPTCHA_TIMEOUT":         unicap.ErrUnsolvable,
	"ERROR_RECAPTCHA_INVALID_SITEKEY": unicap.ErrInvalidTask,
	"ERROR_RECAPTCHA_INVALID_DOMAIN":  unicap.ErrInvalidTask,
	"ERROR_INVALID_KEY_TYPE":          unicap.ErrInvalidTask,
	"ERROR_DOMAIN_NOT_ALLOWED":        unicap.ErrInvalidTask,
	"ERROR_HIGH_DEMAND_DOMAIN":        unicap.ErrInvalidTask,
	"ERROR_FAILED_LOADING_WIDGET":     unicap.ErrInvalidTask,
	"ERROR_ZERO_CAPTCHA_FILESIZE":     unicap.ErrInvalidTask,
	"ERROR_TOO_BIG_CAPTCHA_FILESIZE":  unicap.ErrInvalidTask,
	"ERROR_IMAGE_TYPE_NOT_SUPPORTED":  unicap.ErrInvalidTask,
	"ERROR_ILLEGAL_IMAGE":             unicap.ErrInvalidTask,
	"ERROR_PARSE_IMAGE_FAIL":          unicap.ErrInvalidTask,
	"ERROR_UNKNOWN_QUESTION":          unicap.ErrInvalidTask,
	"ERROR_UNKNOW_450x450_QUESTION":   unicap.ErrInvalidTask,
	"ERROR_PROXY_CONNECT_REFUSED":     unicap.ErrProxy,
	"ERROR_PROXY_CONNECT_TIMEOUT":     unicap.ErrProxy,
	"ERROR_PROXY_READ_TIMEOUT":        unicap.ErrProxy,
	"ERROR_PROXY_BANNED":              unicap.ErrProxy,
	"ERROR_PROXY_INCORRECT":           unicap.ErrProxy,
	"ERROR_PROXY_TRANSPARENT":         unicap.ErrProxy,
	"ERROR_PROXY_NOT_AUTHORISED":      unicap.ErrProxy,
	"ERROR_PROXY_TLS_ERROR":           unicap.ErrProxy,
	"ERROR_PROXY_CONNECTION_FAILED":   unicap.ErrProxy,
	"ERROR_SERVICE_UNAVALIABLE":       unicap.ErrUnavailable, //nolint:misspell // YesCaptcha spells the code this way.
	"ERROR_INTERNAL_SERVER_ERROR":     unicap.ErrUnavailable,
	"ERROR_UNKNOWN":                   unicap.ErrTaskFailed,
}

// mappers registers every task type YesCaptcha solves.
var mappers = solverapi.Mappers{
	unicap.TaskTypeReCaptchaV2:             solverapi.Map(mapReCaptchaV2),
	unicap.TaskTypeReCaptchaV2Enterprise:   solverapi.Map(mapReCaptchaV2Enterprise),
	unicap.TaskTypeReCaptchaV3:             solverapi.Map(mapReCaptchaV3),
	unicap.TaskTypeReCaptchaV3Enterprise:   solverapi.Map(mapReCaptchaV3Enterprise),
	unicap.TaskTypeReCaptchaClassification: solverapi.Map(mapReCaptchaClassification),
	unicap.TaskTypeHCaptcha:                solverapi.Map(mapHCaptcha),
	unicap.TaskTypeTurnstile:               solverapi.Map(mapTurnstile),
	unicap.TaskTypeCloudflareChallenge:     solverapi.Map(mapCloudflareChallenge),
	unicap.TaskTypeImageToText:             solverapi.Map(mapImageToText),
}

// variant returns YesCaptcha's proxyless and proxy wire types for base.
func variant(base string) solverapi.Variant {
	return solverapi.Variant{
		Proxyless: base + "Proxyless",
		Proxy:     base,
	}
}

// proxyless returns the wire type of a task without a proxy variant.
func proxyless(taskType string) solverapi.Variant {
	return solverapi.Variant{
		Proxyless: taskType,
	}
}

// proxyOnly returns the wire type of a task that always needs a proxy.
func proxyOnly(taskType string) solverapi.Variant {
	return solverapi.Variant{
		Proxy: taskType,
	}
}
