package tasks

import (
	"github.com/aarock1234/unicap/v3"
)

// FunCaptchaTask solves an Arkose Labs FunCaptcha challenge.
type FunCaptchaTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// WebsitePublicKey is the Arkose public key, the "pk" parameter.
	WebsitePublicKey string

	// APIJSSubdomain is the custom subdomain api.js loads from, when the
	// site uses one.
	APIJSSubdomain string

	// Data is extra challenge data as a JSON string, such as
	// {"blob":"..."}.
	Data string

	// UserAgent is the browser User-Agent the worker should present. Some
	// providers require it when a proxy is set.
	UserAgent string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeFunCaptcha.
func (t *FunCaptchaTask) Type() unicap.TaskType {
	return unicap.TaskTypeFunCaptcha
}

// Validate reports missing required fields and malformed proxies.
func (t *FunCaptchaTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("WebsitePublicKey", t.WebsitePublicKey),
		optionalProxy(t.Proxy),
	)
}
