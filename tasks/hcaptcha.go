package tasks

import (
	"github.com/aarock1234/unicap/v3"
)

// HCaptchaTask solves an hCaptcha widget.
//
// CapSolver, 2Captcha, and Anti-Captcha have all withdrawn hCaptcha support,
// so the built-in providers reject this task with unicap.ErrUnsupportedTask.
// It remains available for custom providers.
type HCaptchaTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// WebsiteKey is the site key from the widget's data-sitekey attribute.
	WebsiteKey string

	// IsInvisible marks an invisible hCaptcha widget.
	IsInvisible bool

	// EnterprisePayload holds hCaptcha Enterprise parameters such as rqdata.
	// Its keys vary by site, so it is intentionally dynamic.
	EnterprisePayload map[string]any

	// UserAgent is the browser User-Agent the worker should present.
	UserAgent string

	// Cookies are page cookies in "name=value; name2=value2" form.
	Cookies string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeHCaptcha.
func (t *HCaptchaTask) Type() unicap.TaskType {
	return unicap.TaskTypeHCaptcha
}

// Validate reports missing required fields and malformed proxies.
func (t *HCaptchaTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("WebsiteKey", t.WebsiteKey),
		optionalProxy(t.Proxy),
	)
}
