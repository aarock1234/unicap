// Package captchaai provides the CaptchaAI implementation of the unicap
// provider interface, using CaptchaAI's native in.php and res.php form API.
//
// A reply to res.php does not say which kind of answer it holds, so a unicap
// task ID names the shape of the solution as well as the task: it reads
// "<shape>:<task id>", such as "cookie:74965409378".
//
// CaptchaAI sells solver threads instead of credit, so the provider reports
// neither a balance nor a cost. It has no endpoint for reporting an answer,
// and it does not solve hCaptcha.
//
// Proxies are disabled by default on a CaptchaAI account until its support
// enables them. A proxy has to be an HTTP or HTTPS proxy with a username.
//
// Other methods CaptchaAI documents are reachable with a raw task whose
// TaskType is the method and whose Params are the form fields.
package captchaai

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
const Name = "captchaai"

// BaseURL is the API root.
const BaseURL = "https://ocr.captchaai.com"

var (
	_ unicap.Provider       = (*Client)(nil)
	_ unicap.TaskSupporter  = (*Client)(nil)
	_ unicap.PollingAdvisor = (*Client)(nil)
)

// Client is a CaptchaAI provider. It is safe for concurrent use.
type Client struct {
	api *httpapi.Client

	// apiKey travels as a form field of every request, never in a URL.
	apiKey string
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

// New creates a CaptchaAI provider.
func New(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("%s api key is empty: %w", Name, unicap.ErrInvalidAPIKey)
	}

	config := httpapi.Config{
		Name:    Name,
		BaseURL: BaseURL,
	}

	return &Client{
		api:    httpapi.New(config, opts...),
		apiKey: apiKey,
	}, nil
}

// polling follows CaptchaAI's guidance: ask every 5 seconds, after a first
// wait that it puts between 5 and 20 seconds depending on the task type, for
// which 10 seconds is the middle ground. CaptchaAI states no solve deadline,
// so the 3 minute timeout is unicap's choice.
var polling = unicap.Polling{
	InitialDelay:       10 * time.Second,
	Interval:           5 * time.Second,
	MaxInterval:        5 * time.Second,
	Multiplier:         1,
	Timeout:            3 * time.Minute,
	MaxTransientErrors: 3,
}

// errorRules classifies the codes CaptchaAI documents, which arrive in the
// request field of a reply whose status is 0. A nil sentinel marks a
// documented code that matches no unicap condition. The code of a task that
// is not ready is not listed because it is no error. The two server-error
// codes are documented as transient, so they map to unicap.ErrUnavailable;
// no code of a task that failed for good does.
var errorRules = map[string]error{
	"ERROR_WRONG_USER_KEY":           unicap.ErrInvalidAPIKey,
	"ERROR_KEY_DOES_NOT_EXIST":       unicap.ErrInvalidAPIKey,
	"ERROR_ZERO_BALANCE":             unicap.ErrInsufficientFunds, // no free thread, or no active plan
	"ERROR_PAGEURL":                  unicap.ErrInvalidTask,
	"ERROR_TOO_BIG_CAPTCHA_FILESIZE": unicap.ErrInvalidTask,
	"ERROR_ZERO_CAPTCHA_FILESIZE":    unicap.ErrInvalidTask,
	"ERROR_WRONG_FILE_EXTENSION":     unicap.ErrInvalidTask,
	"ERROR_IMAGE_TYPE_NOT_SUPPORTED": unicap.ErrInvalidTask,
	"ERROR_UPLOAD":                   unicap.ErrInvalidTask,
	"IP_BANNED":                      unicap.ErrIPBlocked,
	"ERROR_WRONG_SITEKEY":            unicap.ErrInvalidTask,
	"ERROR_BAD_DOMAIN":               unicap.ErrInvalidTask,
	"ERROR_WRONG_KEY_TYPE":           unicap.ErrInvalidTask,
	"ERROR_BAD_PARAMETERS":           unicap.ErrInvalidTask,
	"ERROR_BAD_PROXY":                unicap.ErrProxy,
	"ERROR_SERVER_ERROR":             unicap.ErrUnavailable,
	"ERROR_INTERNAL_SERVER_ERROR":    unicap.ErrUnavailable,
	"ERROR_CAPTCHA_UNSOLVABLE":       unicap.ErrUnsolvable,
	"ERROR_WRONG_ID_FORMAT":          unicap.ErrTaskNotFound,
	"ERROR_WRONG_CAPTCHA_ID":         unicap.ErrTaskNotFound,
	"ERROR_EMPTY_ACTION":             nil, // a client bug, not a task condition
	"ERROR_PROXY_CONNECTION_FAILED":  unicap.ErrProxy,
}

// mappers registers every task type CaptchaAI solves, each with the shape
// its solution arrives in. Lemin and CaptchaFox are left out: CaptchaAI's
// documentation does not settle the shape of a Lemin solution, and it has no
// field for the User-Agent a CaptchaFox task requires.
var mappers = solverapi.Mappers{
	unicap.TaskTypeReCaptchaV2:           solverapi.Map(shaped(shapeToken, mapReCaptchaV2)),
	unicap.TaskTypeReCaptchaV2Enterprise: solverapi.Map(shaped(shapeAgentToken, mapReCaptchaV2Enterprise)),
	unicap.TaskTypeReCaptchaV3:           solverapi.Map(shaped(shapeToken, mapReCaptchaV3)),
	unicap.TaskTypeReCaptchaV3Enterprise: solverapi.Map(shaped(shapeAgentToken, mapReCaptchaV3Enterprise)),
	unicap.TaskTypeTurnstile:             solverapi.Map(shaped(shapeToken, mapTurnstile)),
	unicap.TaskTypeCloudflareChallenge:   solverapi.Map(shaped(shapeCookie, mapCloudflareChallenge)),
	unicap.TaskTypeGeeTest:               solverapi.Map(shaped(shapeGeeTest, mapGeeTest)),
	unicap.TaskTypeImageToText:           solverapi.Map(shaped(shapeText, mapImageToText)),
	unicap.TaskTypeGrid:                  solverapi.Map(shaped(shapeGrid, mapGrid)),
	unicap.TaskTypeFriendlyCaptcha:       solverapi.Map(shaped(shapeToken, mapFriendlyCaptcha)),
}

// shaped adapts a function that maps a task to its form fields into one that
// returns the whole submission: the fields and the shape of the solution.
func shaped[T unicap.Task](shape solutionShape, mapForm func(T) (form, error)) func(T) (submission, error) {
	return func(task T) (submission, error) {
		fields, err := mapForm(task)
		if err != nil {
			return submission{}, err
		}

		return submission{
			shape:  shape,
			fields: fields,
		}, nil
	}
}
