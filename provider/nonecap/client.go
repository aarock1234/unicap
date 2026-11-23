// Package nonecap provides the NoneCap implementation of the unicap provider
// interface, using NoneCap's native v1 API.
//
// NoneCap describes itself as an open beta and solves hCaptcha only,
// including hCaptcha Enterprise.
//
// NoneCap counts an account's balance in credits, which have no fixed dollar
// value, so the provider reports neither a balance nor a cost.
package nonecap

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
const Name = "nonecap"

// BaseURL is the API root.
const BaseURL = "https://api.nonecap.com"

var (
	_ unicap.Provider       = (*Client)(nil)
	_ unicap.Reporter       = (*Client)(nil)
	_ unicap.TaskSupporter  = (*Client)(nil)
	_ unicap.PollingAdvisor = (*Client)(nil)
)

// Client is a NoneCap provider. It is safe for concurrent use.
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

// New creates a NoneCap provider.
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

// polling follows NoneCap's Go guide, which calls a 2 second ticker a
// sensible default for plain polling. NoneCap states no solve deadline, only
// that most solves finish in well under a minute, so the 3 minute timeout is
// unicap's choice.
var polling = unicap.Polling{
	InitialDelay:       2 * time.Second,
	Interval:           2 * time.Second,
	MaxInterval:        2 * time.Second,
	Multiplier:         1,
	Timeout:            3 * time.Minute,
	MaxTransientErrors: 3,
}

// requestErrorRules classifies the codes of NoneCap's request errors, which
// arrive in an error envelope with an HTTP error status. A code that is not
// listed, such as conflict, unsupported_media_type or expired_window, matches
// no unicap condition and is reported without a sentinel. The forbidden code
// covers a missing scope, an admin-only field and a refused network address
// alike, so it maps to the key being refused.
var requestErrorRules = map[string]error{
	"invalid_request":            unicap.ErrInvalidTask,
	"unauthorized":               unicap.ErrInvalidAPIKey,
	"insufficient_credits":       unicap.ErrInsufficientFunds,
	"key_credit_limit_exceeded":  unicap.ErrInsufficientFunds,
	"account_locked":             unicap.ErrIPBlocked,
	"forbidden":                  unicap.ErrInvalidAPIKey,
	"not_found":                  unicap.ErrTaskNotFound,
	"not_eligible":               unicap.ErrTaskNotFound,
	"payload_too_large":          unicap.ErrInvalidTask,
	"validation_error":           unicap.ErrInvalidTask,
	"recognition_failed":         unicap.ErrUnsolvable,
	"concurrency_limit_exceeded": unicap.ErrUnavailable,
	"proxy_unavailable":          unicap.ErrProxy,
	"rate_capped":                unicap.ErrUnavailable,
	"rate_limited":               unicap.ErrUnavailable,
	"sitekey_rate_limited":       unicap.ErrUnavailable,
	"internal_error":             unicap.ErrUnavailable,
	"maintenance":                unicap.ErrUnavailable,
	"sitekey_unavailable":        unicap.ErrUnavailable,
}

// solveFailureRules classifies the error codes of solves that ended without
// a token. Such a solve is final, so no code maps to unicap.ErrUnavailable,
// which would make the poller ask for the same dead solve again: a failure a
// fresh solve may get past is unicap.ErrUnsolvable instead.
var solveFailureRules = map[string]error{
	"challenge_not_loaded": unicap.ErrUnsolvable,
	"token_not_granted":    unicap.ErrUnsolvable,
	"challenge_expired":    unicap.ErrUnsolvable,
	"challenge_errored":    unicap.ErrUnsolvable,
	"proxy_error":          unicap.ErrProxy,
	"target_unreachable":   unicap.ErrInvalidTask,
	"vision_error":         unicap.ErrUnsolvable,
	"internal_error":       unicap.ErrUnsolvable,
	"unsolvable_variant":   unicap.ErrTaskFailed,
	"capacity_exhausted":   unicap.ErrUnsolvable,
	"cancelled":            unicap.ErrTaskFailed,
	"expired":              unicap.ErrUnsolvable,
}

// solveFailureReasonRules lists the failure reasons that say more than their
// code and so take precedence over solveFailureRules. A stalled proxy expires
// the solve, and a type NoneCap does not serve fails with an internal error
// that no fresh solve gets past.
var solveFailureReasonRules = map[string]error{
	"proxy_stalled":   unicap.ErrProxy,
	"type_not_served": unicap.ErrTaskFailed,
}

// mappers registers every task type NoneCap solves.
var mappers = solverapi.Mappers{
	unicap.TaskTypeHCaptcha: solverapi.Map(mapHCaptcha),
}
