package tasks

import (
	"github.com/aarock1234/unicap/v3"
)

// TurnstileTask solves a Cloudflare Turnstile widget, including the
// Turnstile embedded in Cloudflare challenge pages.
type TurnstileTask struct {
	// WebsiteURL is the page that shows the widget.
	WebsiteURL string

	// WebsiteKey is the site key from the widget's data-sitekey attribute.
	WebsiteKey string

	// Action is the widget's data-action value, when present.
	Action string

	// CData is the widget's data-cdata value, when present.
	CData string

	// PageData is the chlPageData value of a Cloudflare challenge page.
	PageData string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeTurnstile.
func (t *TurnstileTask) Type() unicap.TaskType {
	return unicap.TaskTypeTurnstile
}

// Validate reports missing required fields and malformed proxies.
func (t *TurnstileTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("WebsiteKey", t.WebsiteKey),
		optionalProxy(t.Proxy),
	)
}

// CloudflareChallengeTask passes a Cloudflare "Just a moment" interstitial
// and returns the cf_clearance cookie in unicap.Solution.Cookies. Use the
// same proxy and User-Agent for follow-up requests.
type CloudflareChallengeTask struct {
	// WebsiteURL is the protected page.
	WebsiteURL string

	// HTML is the challenge page body, when available.
	HTML string

	// UserAgent is the Chrome User-Agent the worker should present.
	UserAgent string

	// Proxy is the required proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeCloudflareChallenge.
func (t *CloudflareChallengeTask) Type() unicap.TaskType {
	return unicap.TaskTypeCloudflareChallenge
}

// Validate reports a missing URL or proxy and malformed proxies.
func (t *CloudflareChallengeTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		requiredProxy(t.Proxy),
	)
}
