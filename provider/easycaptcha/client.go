// Package easycaptcha provides the easycaptcha.to implementation of the
// unicap provider interface, using easycaptcha.to's native v1 API.
//
// The provider solves hCaptcha tokens. hCaptcha Enterprise tokens need
// easycaptcha.to to whitelist the account first.
package easycaptcha

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
const Name = "easycaptcha"

// BaseURL is the API root.
const BaseURL = "https://api.easycaptcha.to"

var (
	_ unicap.Provider       = (*Client)(nil)
	_ unicap.BalanceChecker = (*Client)(nil)
	_ unicap.Reporter       = (*Client)(nil)
	_ unicap.TaskSupporter  = (*Client)(nil)
	_ unicap.PollingAdvisor = (*Client)(nil)
)

// Client is an easycaptcha.to provider. It is safe for concurrent use.
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

// New creates an easycaptcha.to provider.
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

// polling is unicap's choice: easycaptcha.to prescribes no interval for the
// native API. Its examples give a task a 300 second deadline, which the
// 5 minute timeout matches.
var polling = unicap.Polling{
	InitialDelay:       2 * time.Second,
	Interval:           2 * time.Second,
	MaxInterval:        2 * time.Second,
	Multiplier:         1,
	Timeout:            5 * time.Minute,
	MaxTransientErrors: 3,
}

// requestErrorRules classifies the codes of easycaptcha.to's request errors,
// which arrive as problem details with an HTTP error status. A nil sentinel
// marks a documented code that matches no unicap condition. The forbidden
// code is one of them because it covers several causes, among them an
// account that is not whitelisted for hCaptcha Enterprise.
var requestErrorRules = map[string]error{
	"invalid_credentials":    unicap.ErrInvalidAPIKey,
	"insufficient_balance":   unicap.ErrInsufficientFunds,
	"task_not_found":         unicap.ErrTaskNotFound,
	"task_expired":           unicap.ErrTaskNotFound,
	"invalid_request":        unicap.ErrInvalidTask,
	"payload_too_large":      unicap.ErrInvalidTask,
	"unsupported_media_type": nil,
	"forbidden":              nil,
	"account_suspended":      unicap.ErrIPBlocked,
	"rate_limited":           unicap.ErrUnavailable,
	"task_failed":            unicap.ErrUnsolvable,
	"internal_error":         unicap.ErrUnavailable,
	"bad_gateway":            unicap.ErrUnavailable,
	"service_unavailable":    unicap.ErrUnavailable,
	"gateway_timeout":        unicap.ErrUnavailable,
	"conflict":               nil,
	"not_found":              nil,
	"method_not_allowed":     nil,
	"extension_outdated":     nil,
}

// taskFailureRules classifies the error codes embedded in failed tasks. A
// failed task is final, so no code maps to unicap.ErrUnavailable, which would
// make the poller ask for the same dead task again. task_failed is the only
// code easycaptcha.to's documentation shows in a failed task.
var taskFailureRules = map[string]error{
	"task_failed": unicap.ErrUnsolvable,
}

// mappers registers every task type easycaptcha.to solves. Its recognition
// types are left out: reCAPTCHA recognition takes image data URIs, which
// unicap's bare base64 cannot name a MIME type for, and no unicap task
// carries an hCaptcha challenge object. Both are reachable with a raw task.
var mappers = solverapi.Mappers{
	unicap.TaskTypeHCaptcha: solverapi.Map(mapHCaptcha),
}
