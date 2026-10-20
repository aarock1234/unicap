// Package twocaptcha provides the 2Captcha implementation of the unicap
// provider interface, using the JSON API v2.
package twocaptcha

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
)

// Name is the provider identifier.
const Name = "2captcha"

// BaseURL is the API root.
const BaseURL = "https://api.2captcha.com"

var (
	_ unicap.Provider       = (*Client)(nil)
	_ unicap.BalanceChecker = (*Client)(nil)
	_ unicap.Reporter       = (*Client)(nil)
	_ unicap.TaskSupporter  = (*Client)(nil)
	_ unicap.PollingAdvisor = (*Client)(nil)
)

// Client is a 2Captcha provider. It is safe for concurrent use.
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

// WithSoftID sets the software ID from the 2Captcha software catalog, which
// credits usage to your application.
func WithSoftID(softID int) Option {
	return solverapi.WithSoftID(softID)
}

// WithCallbackURL sets the pingback URL 2Captcha notifies when a task
// finishes. The URL's domain must be registered in the account settings.
func WithCallbackURL(callbackURL string) Option {
	return solverapi.WithCallbackURL(callbackURL)
}

// WithLanguagePool selects the worker language pool for image and text
// captchas, such as "en" or "ru".
func WithLanguagePool(languagePool string) Option {
	return solverapi.WithLanguagePool(languagePool)
}

// New creates a 2Captcha provider.
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

// polling follows 2Captcha's guidance to poll every 5 seconds; faster
// polling can get the account temporarily blocked.
var polling = unicap.Polling{
	InitialDelay:       5 * time.Second,
	Interval:           5 * time.Second,
	MaxInterval:        5 * time.Second,
	Multiplier:         1,
	Timeout:            5 * time.Minute,
	MaxTransientErrors: 3,
}

// errorRules classifies 2Captcha's documented error codes.
var errorRules = solverapi.ErrorRules{
	"ERROR_KEY_DOES_NOT_EXIST":        unicap.ErrInvalidAPIKey,
	"ERROR_ZERO_BALANCE":              unicap.ErrInsufficientFunds,
	"ERROR_NO_SLOT_AVAILABLE":         unicap.ErrUnavailable,
	"ERROR_NO_SUCH_CAPCHA_ID":         unicap.ErrTaskNotFound,
	"ERROR_ZERO_CAPTCHA_FILESIZE":     unicap.ErrInvalidTask,
	"ERROR_TOO_BIG_CAPTCHA_FILESIZE":  unicap.ErrInvalidTask,
	"ERROR_PAGEURL":                   unicap.ErrInvalidTask,
	"ERROR_IMAGE_TYPE_NOT_SUPPORTED":  unicap.ErrInvalidTask,
	"ERROR_TASK_ABSENT":               unicap.ErrInvalidTask,
	"ERROR_RECAPTCHA_INVALID_SITEKEY": unicap.ErrInvalidTask,
	"ERROR_BAD_PARAMETERS":            unicap.ErrInvalidTask,
	"ERROR_BAD_IMGINSTRUCTIONS":       unicap.ErrInvalidTask,
	"ERROR_TASK_NOT_SUPPORTED":        unicap.ErrUnsupportedTask,
	"ERROR_CAPTCHA_UNSOLVABLE":        unicap.ErrUnsolvable,
	"ERROR_BAD_DUPLICATES":            unicap.ErrUnsolvable,
	"ERROR_BAD_PROXY":                 unicap.ErrProxy,
	"ERROR_IP_NOT_ALLOWED":            unicap.ErrIPBlocked,
	"ERROR_IP_BLOCKED":                unicap.ErrIPBlocked,
	"ERROR_ACCOUNT_SUSPENDED":         unicap.ErrIPBlocked,
}

// mappers registers every task type 2Captcha solves.
var mappers = solverapi.Mappers{
	unicap.TaskTypeReCaptchaV2:           solverapi.Map(mapReCaptchaV2),
	unicap.TaskTypeReCaptchaV2Enterprise: solverapi.Map(mapReCaptchaV2Enterprise),
	unicap.TaskTypeReCaptchaV3:           solverapi.Map(mapReCaptchaV3),
	unicap.TaskTypeReCaptchaV3Enterprise: solverapi.Map(mapReCaptchaV3Enterprise),
	unicap.TaskTypeFunCaptcha:            solverapi.Map(mapFunCaptcha),
	unicap.TaskTypeTurnstile:             solverapi.Map(mapTurnstile),
	unicap.TaskTypeDataDome:              solverapi.Map(mapDataDome),
	unicap.TaskTypeGeeTest:               solverapi.Map(mapGeeTest),
	unicap.TaskTypeGeeTestV4:             solverapi.Map(mapGeeTestV4),
	unicap.TaskTypeAWSWAF:                solverapi.Map(mapAWSWAF),
	unicap.TaskTypeMTCaptcha:             solverapi.Map(mapMTCaptcha),
	unicap.TaskTypeFriendlyCaptcha:       solverapi.Map(mapFriendlyCaptcha),
	unicap.TaskTypeLemin:                 solverapi.Map(mapLemin),
	unicap.TaskTypeCutCaptcha:            solverapi.Map(mapCutCaptcha),
	unicap.TaskTypeProsopo:               solverapi.Map(mapProsopo),
	unicap.TaskTypeAltcha:                solverapi.Map(mapAltcha),
	unicap.TaskTypeCapy:                  solverapi.Map(mapCapy),
	unicap.TaskTypeKeyCaptcha:            solverapi.Map(mapKeyCaptcha),
	unicap.TaskTypeCyberSiAra:            solverapi.Map(mapCyberSiAra),
	unicap.TaskTypeAtbCaptcha:            solverapi.Map(mapAtbCaptcha),
	unicap.TaskTypeTencent:               solverapi.Map(mapTencent),
	unicap.TaskTypeCaptchaFox:            solverapi.Map(mapCaptchaFox),
	unicap.TaskTypeYandexSmartCaptcha:    solverapi.Map(mapYandexSmartCaptcha),
	unicap.TaskTypeVKCaptcha:             solverapi.Map(mapVKCaptcha),
	unicap.TaskTypeYidun:                 solverapi.Map(mapYidun),
	unicap.TaskTypeBinance:               solverapi.Map(mapBinance),
	unicap.TaskTypeHunt:                  solverapi.Map(mapHunt),
	unicap.TaskTypeTSPD:                  solverapi.Map(mapTSPD),
	unicap.TaskTypeBasilisk:              solverapi.Map(mapBasilisk),
	unicap.TaskTypeImperva:               solverapi.Map(mapImperva),
	unicap.TaskTypeAlibaba:               solverapi.Map(mapAlibaba),
	unicap.TaskTypeImageToText:           solverapi.Map(mapImageToText),
	unicap.TaskTypeText:                  solverapi.Map(mapText),
	unicap.TaskTypeAudio:                 solverapi.Map(mapAudio),
	unicap.TaskTypeCoordinates:           solverapi.Map(mapCoordinates),
	unicap.TaskTypeGrid:                  solverapi.Map(mapGrid),
	unicap.TaskTypeRotate:                solverapi.Map(mapRotate),
	unicap.TaskTypeDrawAround:            solverapi.Map(mapDrawAround),
	unicap.TaskTypeBoundingBox:           solverapi.Map(mapBoundingBox),
	unicap.TaskTypeDragAndDrop:           solverapi.Map(mapDragAndDrop),
	unicap.TaskTypeTemu:                  solverapi.Map(mapTemu),
	unicap.TaskTypeVKCaptchaImage:        solverapi.Map(mapVKCaptchaImage),
	unicap.TaskTypeYandexImage:           solverapi.Map(mapYandexImage),
	unicap.TaskTypeYandexPuzzle:          solverapi.Map(mapYandexPuzzle),
}

// variant returns 2Captcha's proxyless and proxy wire types for base.
func variant(base string) solverapi.Variant {
	return solverapi.Variant{
		Proxyless: base + "Proxyless",
		Proxy:     base,
	}
}

// proxyOnly returns the wire type of a task that always needs a proxy.
func proxyOnly(taskType string) solverapi.Variant {
	return solverapi.Variant{
		Proxy: taskType,
	}
}

// proxyless returns the wire type of a task without a proxy variant.
func proxyless(taskType string) solverapi.Variant {
	return solverapi.Variant{
		Proxyless: taskType,
	}
}
