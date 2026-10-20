package tasks

import (
	"github.com/aarock1234/unicap/v3"
)

// ReCaptchaV2Task solves a Google reCAPTCHA v2 checkbox or invisible widget.
type ReCaptchaV2Task struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// WebsiteKey is the site key from the widget's data-sitekey attribute.
	WebsiteKey string

	// IsInvisible marks an invisible reCAPTCHA widget.
	IsInvisible bool

	// PageAction is the "sa" value of the anchor request, when present.
	PageAction string

	// DataS is the data-s value some Google services set on the widget.
	DataS string

	// APIDomain is the domain the widget script loads from, either
	// "www.google.com" or "www.recaptcha.net".
	APIDomain string

	// IsSession asks the provider to return session cookies with the token.
	IsSession bool

	// UserAgent is the browser User-Agent the worker should present.
	UserAgent string

	// Cookies are page cookies in "name=value; name2=value2" form.
	Cookies string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeReCaptchaV2.
func (t *ReCaptchaV2Task) Type() unicap.TaskType {
	return unicap.TaskTypeReCaptchaV2
}

// Validate reports missing required fields and malformed proxies.
func (t *ReCaptchaV2Task) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("WebsiteKey", t.WebsiteKey),
		optionalProxy(t.Proxy),
	)
}

// ReCaptchaV2EnterpriseTask solves a reCAPTCHA v2 Enterprise widget.
type ReCaptchaV2EnterpriseTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// WebsiteKey is the site key from the widget's data-sitekey attribute.
	WebsiteKey string

	// IsInvisible marks an invisible reCAPTCHA widget.
	IsInvisible bool

	// PageAction is the "sa" value of the anchor request, when present.
	PageAction string

	// EnterprisePayload holds extra grecaptcha.enterprise.render parameters,
	// such as {"s": "..."}. Its keys are defined by Google and vary by site,
	// so it is intentionally dynamic.
	EnterprisePayload map[string]any

	// APIDomain is the domain the widget script loads from, either
	// "www.google.com" or "www.recaptcha.net".
	APIDomain string

	// IsSession asks the provider to return session cookies with the token.
	IsSession bool

	// UserAgent is the browser User-Agent the worker should present.
	UserAgent string

	// Cookies are page cookies in "name=value; name2=value2" form.
	Cookies string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeReCaptchaV2Enterprise.
func (t *ReCaptchaV2EnterpriseTask) Type() unicap.TaskType {
	return unicap.TaskTypeReCaptchaV2Enterprise
}

// Validate reports missing required fields and malformed proxies.
func (t *ReCaptchaV2EnterpriseTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("WebsiteKey", t.WebsiteKey),
		optionalProxy(t.Proxy),
	)
}

// ReCaptchaV3Task solves a score-based reCAPTCHA v3 challenge.
type ReCaptchaV3Task struct {
	// WebsiteURL is the page that runs the captcha.
	WebsiteURL string

	// WebsiteKey is the site key passed to grecaptcha.execute.
	WebsiteKey string

	// PageAction is the action passed to grecaptcha.execute.
	PageAction string

	// MinScore is the minimum score the token should achieve, between 0 and
	// 1. Zero leaves the choice to the provider.
	MinScore float64

	// HighScore requests a provider's premium high-score solving mode.
	HighScore bool

	// APIDomain is the domain the script loads from, either
	// "www.google.com" or "www.recaptcha.net".
	APIDomain string

	// IsSession asks the provider to return session cookies with the token.
	IsSession bool

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeReCaptchaV3.
func (t *ReCaptchaV3Task) Type() unicap.TaskType {
	return unicap.TaskTypeReCaptchaV3
}

// Validate reports missing required fields, an out-of-range MinScore, and
// malformed proxies.
func (t *ReCaptchaV3Task) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("WebsiteKey", t.WebsiteKey),
		inRange("MinScore", t.MinScore, 0, 1),
		optionalProxy(t.Proxy),
	)
}

// ReCaptchaV3EnterpriseTask solves a score-based reCAPTCHA v3 Enterprise
// challenge.
type ReCaptchaV3EnterpriseTask struct {
	// WebsiteURL is the page that runs the captcha.
	WebsiteURL string

	// WebsiteKey is the site key passed to grecaptcha.enterprise.execute.
	WebsiteKey string

	// PageAction is the action passed to grecaptcha.enterprise.execute.
	PageAction string

	// MinScore is the minimum score the token should achieve, between 0 and
	// 1. Zero leaves the choice to the provider.
	MinScore float64

	// HighScore requests a provider's premium high-score solving mode.
	HighScore bool

	// EnterprisePayload holds extra enterprise parameters, such as
	// {"s": "..."}. Its keys vary by site, so it is intentionally dynamic.
	EnterprisePayload map[string]any

	// APIDomain is the domain the script loads from, either
	// "www.google.com" or "www.recaptcha.net".
	APIDomain string

	// IsSession asks the provider to return session cookies with the token.
	IsSession bool

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeReCaptchaV3Enterprise.
func (t *ReCaptchaV3EnterpriseTask) Type() unicap.TaskType {
	return unicap.TaskTypeReCaptchaV3Enterprise
}

// Validate reports missing required fields, an out-of-range MinScore, and
// malformed proxies.
func (t *ReCaptchaV3EnterpriseTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("WebsiteKey", t.WebsiteKey),
		inRange("MinScore", t.MinScore, 0, 1),
		optionalProxy(t.Proxy),
	)
}

// ReCaptchaClassificationTask identifies the tiles of a reCAPTCHA image grid
// that contain the requested object.
type ReCaptchaClassificationTask struct {
	// Image is the base64-encoded grid image.
	Image string

	// Question is the object to find as a Freebase ID, such as "/m/04_sv"
	// for motorcycles.
	Question string

	// WebsiteURL is the optional page that shows the captcha.
	WebsiteURL string

	// WebsiteKey is the optional reCAPTCHA site key.
	WebsiteKey string
}

// Type returns unicap.TaskTypeReCaptchaClassification.
func (t *ReCaptchaClassificationTask) Type() unicap.TaskType {
	return unicap.TaskTypeReCaptchaClassification
}

// Validate reports missing required fields.
func (t *ReCaptchaClassificationTask) Validate() error {
	return validate(
		required("Image", t.Image),
		required("Question", t.Question),
	)
}

// ReCaptchaClassificationSolution is the decoded payload of a
// ReCaptchaClassificationTask. Matching tiles are also available in
// unicap.Solution.Indexes.
type ReCaptchaClassificationSolution struct {
	// Type is "multi" for 3x3 or 4x4 grids and "single" for one image.
	Type string `json:"type"`

	// Objects are the zero-based indexes of matching tiles.
	Objects []int `json:"objects"`

	// HasObject reports whether a single image contains the object.
	HasObject bool `json:"hasObject"`

	// Size is the grid size.
	Size int `json:"size"`
}
