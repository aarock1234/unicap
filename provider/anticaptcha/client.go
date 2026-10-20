// Package anticaptcha provides the Anti-Captcha implementation of the unicap
// provider interface.
//
// Anti-Captcha accepts proxies by IP address only and requires a User-Agent
// for reCAPTCHA v2, FunCaptcha, and GeeTest tasks that use a proxy.
package anticaptcha

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
)

// Name is the provider identifier.
const Name = "anticaptcha"

// BaseURL is the API root.
const BaseURL = "https://api.anti-captcha.com"

var (
	_ unicap.Provider       = (*Client)(nil)
	_ unicap.BalanceChecker = (*Client)(nil)
	_ unicap.Reporter       = (*Client)(nil)
	_ unicap.TaskSupporter  = (*Client)(nil)
	_ unicap.PollingAdvisor = (*Client)(nil)
)

// Client is an Anti-Captcha provider. It is safe for concurrent use.
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

// WithSoftID sets the application ID from the Anti-Captcha developer center,
// which credits usage to your application.
func WithSoftID(softID int) Option {
	return solverapi.WithSoftID(softID)
}

// WithCallbackURL sets the URL Anti-Captcha posts results to when a task
// finishes.
func WithCallbackURL(callbackURL string) Option {
	return solverapi.WithCallbackURL(callbackURL)
}

// New creates an Anti-Captcha provider.
func New(apiKey string, opts ...Option) (*Client, error) {
	config := solverapi.Config{
		Name:         Name,
		BaseURL:      BaseURL,
		APIKey:       apiKey,
		TaskIDFormat: solverapi.TaskIDNumber,
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

// polling follows Anti-Captcha's guidance: wait about 5 seconds, then poll
// every few seconds. Finished results are kept for only 60 seconds.
var polling = unicap.Polling{
	InitialDelay:       5 * time.Second,
	Interval:           3 * time.Second,
	MaxInterval:        3 * time.Second,
	Multiplier:         1,
	Timeout:            5 * time.Minute,
	MaxTransientErrors: 3,
}

// errorRules classifies Anti-Captcha's documented error codes.
var errorRules = solverapi.ErrorRules{
	"ERROR_KEY_DOES_NOT_EXIST":              unicap.ErrInvalidAPIKey,
	"ERROR_ZERO_BALANCE":                    unicap.ErrInsufficientFunds,
	"ERROR_NO_SLOT_AVAILABLE":               unicap.ErrUnavailable,
	"ERROR_NO_SUCH_CAPCHA_ID":               unicap.ErrTaskNotFound,
	"ERROR_ZERO_CAPTCHA_FILESIZE":           unicap.ErrInvalidTask,
	"ERROR_TOO_BIG_CAPTCHA_FILESIZE":        unicap.ErrInvalidTask,
	"ERROR_IMAGE_TYPE_NOT_SUPPORTED":        unicap.ErrInvalidTask,
	"ERROR_TASK_ABSENT":                     unicap.ErrInvalidTask,
	"ERROR_INCORRECT_SESSION_DATA":          unicap.ErrInvalidTask,
	"ERROR_RECAPTCHA_INVALID_SITEKEY":       unicap.ErrInvalidTask,
	"ERROR_RECAPTCHA_INVALID_DOMAIN":        unicap.ErrInvalidTask,
	"ERROR_RECAPTCHA_OLD_BROWSER":           unicap.ErrInvalidTask,
	"ERROR_TOKEN_EXPIRED":                   unicap.ErrInvalidTask,
	"ERROR_INVALID_KEY_TYPE":                unicap.ErrInvalidTask,
	"ERROR_VISIBLE_RECAPTCHA":               unicap.ErrInvalidTask,
	"ERROR_TEMPLATE_NOT_FOUND":              unicap.ErrInvalidTask,
	"ERROR_TASK_NOT_SUPPORTED":              unicap.ErrUnsupportedTask,
	"ERROR_CAPTCHA_UNSOLVABLE":              unicap.ErrUnsolvable,
	"ERROR_BAD_DUPLICATES":                  unicap.ErrUnsolvable,
	"ERROR_RECAPTCHA_TIMEOUT":               unicap.ErrUnsolvable,
	"ERROR_FAILED_LOADING_WIDGET":           unicap.ErrUnsolvable,
	"ERROR_TASK_CANCELED":                   unicap.ErrTaskFailed,
	"ERROR_PROXY_CONNECT_REFUSED":           unicap.ErrProxy,
	"ERROR_PROXY_CONNECT_TIMEOUT":           unicap.ErrProxy,
	"ERROR_PROXY_READ_TIMEOUT":              unicap.ErrProxy,
	"ERROR_PROXY_BANNED":                    unicap.ErrProxy,
	"ERROR_PROXY_TRANSPARENT":               unicap.ErrProxy,
	"ERROR_PROXY_HAS_NO_IMAGE_SUPPORT":      unicap.ErrProxy,
	"ERROR_PROXY_INCOMPATIBLE_HTTP_VERSION": unicap.ErrProxy,
	"ERROR_PROXY_NOT_AUTHORISED":            unicap.ErrProxy,
	"ERROR_IP_NOT_ALLOWED":                  unicap.ErrIPBlocked,
	"ERROR_IP_BLOCKED":                      unicap.ErrIPBlocked,
	"ERROR_ACCOUNT_SUSPENDED":               unicap.ErrIPBlocked,
}

// mappers registers every task type Anti-Captcha solves.
var mappers = solverapi.Mappers{
	unicap.TaskTypeReCaptchaV2:           solverapi.Map(mapReCaptchaV2),
	unicap.TaskTypeReCaptchaV2Enterprise: solverapi.Map(mapReCaptchaV2Enterprise),
	unicap.TaskTypeReCaptchaV3:           solverapi.Map(mapReCaptchaV3),
	unicap.TaskTypeReCaptchaV3Enterprise: solverapi.Map(mapReCaptchaV3Enterprise),
	unicap.TaskTypeFunCaptcha:            solverapi.Map(mapFunCaptcha),
	unicap.TaskTypeTurnstile:             solverapi.Map(mapTurnstile),
	unicap.TaskTypeGeeTest:               solverapi.Map(mapGeeTest),
	unicap.TaskTypeGeeTestV4:             solverapi.Map(mapGeeTestV4),
	unicap.TaskTypeAWSWAF:                solverapi.Map(mapAWSWAF),
	unicap.TaskTypeFriendlyCaptcha:       solverapi.Map(mapFriendlyCaptcha),
	unicap.TaskTypeProsopo:               solverapi.Map(mapProsopo),
	unicap.TaskTypeAltcha:                solverapi.Map(mapAltcha),
	unicap.TaskTypeAntiGate:              solverapi.Map(mapAntiGate),
	unicap.TaskTypeImageToText:           solverapi.Map(mapImageToText),
	unicap.TaskTypeCoordinates:           solverapi.Map(mapCoordinates),
}

// variant returns Anti-Captcha's proxyless and proxy wire types for base.
func variant(base string) solverapi.Variant {
	return solverapi.Variant{
		Proxyless: base + "Proxyless",
		Proxy:     base,
	}
}
