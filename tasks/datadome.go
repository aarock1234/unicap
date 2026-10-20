package tasks

import (
	"github.com/aarock1234/unicap/v3"
)

// DataDomeTask solves a DataDome slider and returns the datadome cookie in
// unicap.Solution.Cookies. Use the same proxy and User-Agent for follow-up
// requests.
type DataDomeTask struct {
	// WebsiteURL is the protected page.
	WebsiteURL string

	// CaptchaURL is the captcha iframe URL. Its "t" parameter must be "fe";
	// "t=bv" means the proxy IP is banned.
	CaptchaURL string

	// UserAgent is the browser User-Agent the worker should present.
	UserAgent string

	// HTML is the blocked page body, when available.
	HTML string

	// Proxy is the required proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeDataDome.
func (t *DataDomeTask) Type() unicap.TaskType {
	return unicap.TaskTypeDataDome
}

// Validate reports missing required fields and a missing or malformed proxy.
func (t *DataDomeTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("CaptchaURL", t.CaptchaURL),
		required("UserAgent", t.UserAgent),
		requiredProxy(t.Proxy),
	)
}
