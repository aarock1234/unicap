// Package nopecha provides the NopeCHA implementation of the unicap provider
// interface, using NopeCHA's native v1 API.
//
// NopeCHA fetches a job's result from the endpoint the job was submitted to,
// so a unicap task ID names both: it reads "<endpoint>:<job id>", such as
// "token/hcaptcha:2xut...".
//
// The provider solves hCaptcha, reCAPTCHA and Turnstile tokens. NopeCHA's
// recognition endpoints are reachable with a raw task whose TaskType is the
// endpoint, such as "recognition/funcaptcha".
//
// NopeCHA counts an account's balance in credits, which have no fixed dollar
// value, so the provider reports neither a balance nor a cost.
package nopecha

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
const Name = "nopecha"

// BaseURL is the API root.
const BaseURL = "https://api.nopecha.com"

var (
	_ unicap.Provider       = (*Client)(nil)
	_ unicap.TaskSupporter  = (*Client)(nil)
	_ unicap.PollingAdvisor = (*Client)(nil)
)

// Client is a NopeCHA provider. It is safe for concurrent use.
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

// New creates a NopeCHA provider.
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

// polling follows NopeCHA's advice to wait before asking again for a job
// that is not finished. NopeCHA recommends 500 milliseconds; a fixed
// 1 second pace keeps well under its rate limit, which applies to each key
// on each endpoint. NopeCHA states no solve deadline, so the 3 minute
// timeout is unicap's choice.
var polling = unicap.Polling{
	InitialDelay:       time.Second,
	Interval:           time.Second,
	MaxInterval:        time.Second,
	Multiplier:         1,
	Timeout:            3 * time.Minute,
	MaxTransientErrors: 3,
}

// errorRules classifies NopeCHA's error codes, which arrive as integers in
// an error body with an HTTP error status. A nil sentinel marks a documented
// code that matches no unicap condition. The incomplete-job code is not
// listed because it is no error when a job's result is fetched. NopeCHA
// documents no code for a job that failed, expired or is unknown.
var errorRules = map[int]error{
	9:  unicap.ErrUnavailable,       // unknown error, which NopeCHA asks to retry
	10: unicap.ErrInvalidTask,       // invalid request
	11: unicap.ErrUnavailable,       // rate limit reached
	12: unicap.ErrIPBlocked,         // banned ip
	15: unicap.ErrInvalidAPIKey,     // invalid key
	16: unicap.ErrInsufficientFunds, // out of credit
	17: nil,                         // update required
	18: unicap.ErrUnsupportedTask,   // feature unavailable for the current plan
	19: unicap.ErrUnsupportedTask,   // endpoint discontinued
}

// mappers registers every task type NopeCHA solves, each with the endpoint
// that takes it. NopeCHA's recognition endpoints are left out: their
// documentation shows images as data URIs, which unicap's bare base64 cannot
// name a MIME type for, and several match no unicap task. They are reachable
// with a raw task.
var mappers = solverapi.Mappers{
	unicap.TaskTypeHCaptcha:              solverapi.Map(at(endpointHCaptcha, mapHCaptcha)),
	unicap.TaskTypeReCaptchaV2:           solverapi.Map(at(endpointReCaptchaV2, mapReCaptchaV2)),
	unicap.TaskTypeReCaptchaV2Enterprise: solverapi.Map(at(endpointReCaptchaV2, mapReCaptchaV2Enterprise)),
	unicap.TaskTypeReCaptchaV3:           solverapi.Map(at(endpointReCaptchaV3, mapReCaptchaV3)),
	unicap.TaskTypeReCaptchaV3Enterprise: solverapi.Map(at(endpointReCaptchaV3, mapReCaptchaV3Enterprise)),
	unicap.TaskTypeTurnstile:             solverapi.Map(at(endpointTurnstile, mapTurnstile)),
}

// at adapts a function that maps a task to the body of one endpoint into one
// that returns the whole submission to that endpoint.
func at[T unicap.Task, B any](endpoint string, mapBody func(T) (B, error)) func(T) (endpointBody[B], error) {
	return func(task T) (endpointBody[B], error) {
		body, err := mapBody(task)
		if err != nil {
			return endpointBody[B]{}, err
		}

		return endpointBody[B]{
			endpoint: endpoint,
			body:     body,
		}, nil
	}
}
