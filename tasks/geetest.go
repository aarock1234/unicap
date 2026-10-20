package tasks

import (
	"encoding/json"

	"github.com/aarock1234/unicap/v3"
)

// GeeTestTask solves a GeeTest v3 challenge. Decode the result with
// unicap.DecodeSolution[GeeTestSolution].
type GeeTestTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// GT is the static gt value from the captcha configuration.
	GT string

	// Challenge is the fresh challenge value. It changes on every page load.
	Challenge string

	// APIServerSubdomain is the custom GeeTest API host, when the site uses
	// one.
	APIServerSubdomain string

	// UserAgent is the browser User-Agent the worker should present. Some
	// providers require it when a proxy is set.
	UserAgent string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeGeeTest.
func (t *GeeTestTask) Type() unicap.TaskType {
	return unicap.TaskTypeGeeTest
}

// Validate reports missing required fields and malformed proxies.
func (t *GeeTestTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("GT", t.GT),
		required("Challenge", t.Challenge),
		optionalProxy(t.Proxy),
	)
}

// GeeTestSolution is the decoded payload of a GeeTestTask.
type GeeTestSolution struct {
	Challenge string `json:"challenge"`
	Validate  string `json:"validate"`
	Seccode   string `json:"seccode"`
}

// GeeTestV4Task solves a GeeTest v4 challenge. Decode the result with
// unicap.DecodeSolution[GeeTestV4Solution].
type GeeTestV4Task struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// CaptchaID is the captcha_id from the captcha configuration.
	CaptchaID string

	// RiskType is the optional risk_type the site initializes the captcha
	// with.
	RiskType string

	// APIServerSubdomain is the custom GeeTest API host, when the site uses
	// one.
	APIServerSubdomain string

	// UserAgent is the browser User-Agent the worker should present. Some
	// providers require it when a proxy is set.
	UserAgent string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeGeeTestV4.
func (t *GeeTestV4Task) Type() unicap.TaskType {
	return unicap.TaskTypeGeeTestV4
}

// Validate reports missing required fields and malformed proxies.
func (t *GeeTestV4Task) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("CaptchaID", t.CaptchaID),
		optionalProxy(t.Proxy),
	)
}

// GeeTestV4Solution is the decoded payload of a GeeTestV4Task.
type GeeTestV4Solution struct {
	CaptchaID     string      `json:"captcha_id"`
	LotNumber     string      `json:"lot_number"`
	PassToken     string      `json:"pass_token"`
	GenTime       json.Number `json:"gen_time"`
	CaptchaOutput string      `json:"captcha_output"`
}
