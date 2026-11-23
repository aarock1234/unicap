// Package capmonster provides the CapMonster Cloud implementation of the
// unicap provider interface, using CapMonster Cloud's native API.
//
// Every CapMonster Cloud task is asynchronous: CreateTask returns a pending
// task and the solution comes from GetTaskResult.
//
// CapMonster Cloud does not offer hCaptcha.
package capmonster

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
)

// Name is the provider identifier.
const Name = "capmonster"

// BaseURL is the API root.
const BaseURL = "https://api.capmonster.cloud"

var (
	_ unicap.Provider       = (*Client)(nil)
	_ unicap.BalanceChecker = (*Client)(nil)
	_ unicap.Reporter       = (*Client)(nil)
	_ unicap.TaskSupporter  = (*Client)(nil)
	_ unicap.PollingAdvisor = (*Client)(nil)
)

// Client is a CapMonster Cloud provider. It is safe for concurrent use.
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

// WithCallbackURL sets the URL CapMonster Cloud posts results to when a task
// finishes. The endpoint must accept the request within 2 seconds.
func WithCallbackURL(callbackURL string) Option {
	return solverapi.WithCallbackURL(callbackURL)
}

// New creates a CapMonster Cloud provider.
func New(apiKey string, opts ...Option) (*Client, error) {
	config := solverapi.Config{
		Name:         Name,
		BaseURL:      BaseURL,
		APIKey:       apiKey,
		TaskIDFormat: solverapi.TaskIDNumber,
		Mappers:      mappers,
		Errors:       errorRules,
		Polling:      polling,
		Extras: []solverapi.Extra{
			solverapi.ExtraCallbackURL,
		},
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

// polling follows CapMonster Cloud's limits: a task may be polled no more
// often than once every 2 seconds and at most 120 times, or the account may
// be temporarily locked. A fixed 3 second pace polls 100 times in the 5
// minute timeout.
var polling = unicap.Polling{
	InitialDelay:       3 * time.Second,
	Interval:           3 * time.Second,
	MaxInterval:        3 * time.Second,
	Multiplier:         1,
	Timeout:            5 * time.Minute,
	MaxTransientErrors: 3,
}

// errorRules classifies CapMonster Cloud's documented error codes.
// CAPTCHA_NOT_READY means the task is still being solved, so it maps to the
// one sentinel the poller treats as transient.
var errorRules = solverapi.ErrorRules{
	"ERROR_KEY_DOES_NOT_EXIST":                  unicap.ErrInvalidAPIKey,
	"ERROR_ZERO_BALANCE":                        unicap.ErrInsufficientFunds,
	"ERROR_NO_SUCH_CAPCHA_ID":                   unicap.ErrTaskNotFound,
	"WRONG_CAPTCHA_ID":                          unicap.ErrTaskNotFound,
	"ERROR_CAPTCHA_UNSOLVABLE":                  unicap.ErrUnsolvable,
	"ERROR_RECAPTCHA_TIMEOUT":                   unicap.ErrUnsolvable,
	"CAPTCHA_NOT_READY":                         unicap.ErrUnavailable,
	"ERROR_TOO_MUCH_REQUESTS":                   unicap.ErrUnavailable,
	"ERROR_SERVICE_NOT_AVAILABLE":               unicap.ErrUnavailable,
	"ERROR_IP_NOT_ALLOWED":                      unicap.ErrIPBlocked,
	"ERROR_IP_BANNED":                           unicap.ErrIPBlocked,
	"ERROR_IP_BLOCKED":                          unicap.ErrIPBlocked,
	"ERROR_NO_SUCH_METHOD":                      unicap.ErrUnsupportedTask,
	"ERROR_TASK_NOT_SUPPORTED":                  unicap.ErrUnsupportedTask,
	"ERROR_DOMAIN_NOT_ALLOWED":                  unicap.ErrInvalidTask,
	"ERROR_TASK_ABSENT":                         unicap.ErrInvalidTask,
	"ERROR_INVALID_TASK":                        unicap.ErrInvalidTask,
	"ERROR_TOKEN_EXPIRED":                       unicap.ErrInvalidTask,
	"ERROR_RECAPTCHA_INVALID_SITEKEY":           unicap.ErrInvalidTask,
	"ERROR_RECAPTCHA_INVALID_DOMAIN":            unicap.ErrInvalidTask,
	"ERROR_WRONG_USERAGENT":                     unicap.ErrInvalidTask,
	"ERROR_TOO_BIG_CAPTCHA_FILESIZE":            unicap.ErrInvalidTask,
	"ERROR_ZERO_CAPTCHA_FILESIZE":               unicap.ErrInvalidTask,
	"ERROR_PROXY_CREDENTIALS_INVALID_CHARACTER": unicap.ErrProxy,
	"ERROR_PROXY_CONNECT_REFUSED":               unicap.ErrProxy,
	"ERROR_PROXY_BANNED":                        unicap.ErrProxy,
	"ERROR_PROXY_MISSING":                       unicap.ErrProxy,
	"ERROR_PROXY_NOT_AUTHORISED":                unicap.ErrProxy,
	"ERROR_PROXY_READ_TIMEOUT":                  unicap.ErrProxy,
}

// mappers registers every task type CapMonster Cloud solves.
var mappers = solverapi.Mappers{
	unicap.TaskTypeReCaptchaV2:           solverapi.Map(mapReCaptchaV2),
	unicap.TaskTypeReCaptchaV2Enterprise: solverapi.Map(mapReCaptchaV2Enterprise),
	unicap.TaskTypeReCaptchaV3:           solverapi.Map(mapReCaptchaV3),
	unicap.TaskTypeReCaptchaV3Enterprise: solverapi.Map(mapReCaptchaV3Enterprise),
	unicap.TaskTypeFunCaptcha:            solverapi.Map(mapFunCaptcha),
	unicap.TaskTypeTurnstile:             solverapi.Map(mapTurnstile),
	unicap.TaskTypeCloudflareChallenge:   solverapi.Map(mapCloudflareChallenge),
	unicap.TaskTypeGeeTest:               solverapi.Map(mapGeeTest),
	unicap.TaskTypeGeeTestV4:             solverapi.Map(mapGeeTestV4),
	unicap.TaskTypeAWSWAF:                solverapi.Map(mapAWSWAF),
	unicap.TaskTypeMTCaptcha:             solverapi.Map(mapMTCaptcha),
	unicap.TaskTypeProsopo:               solverapi.Map(mapProsopo),
	unicap.TaskTypeYidun:                 solverapi.Map(mapYidun),
	unicap.TaskTypeBinance:               solverapi.Map(mapBinance),
	unicap.TaskTypeFriendlyCaptcha:       solverapi.Map(mapFriendlyCaptcha),
	unicap.TaskTypeTencent:               solverapi.Map(mapTencent),
	unicap.TaskTypeHunt:                  solverapi.Map(mapHunt),
	unicap.TaskTypeTSPD:                  solverapi.Map(mapTSPD),
	unicap.TaskTypeBasilisk:              solverapi.Map(mapBasilisk),
	unicap.TaskTypeImperva:               solverapi.Map(mapImperva),
	unicap.TaskTypeAlibaba:               solverapi.Map(mapAlibaba),
	unicap.TaskTypeImageToText:           solverapi.Map(mapImageToText),
}
