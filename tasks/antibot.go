package tasks

import (
	"github.com/aarock1234/unicap/v3"
)

// VKCaptchaTask solves a VK ID captcha and returns a token.
type VKCaptchaTask struct {
	// RedirectURI is the captcha redirect_uri from the VK API error.
	RedirectURI string

	// UserAgent is the browser User-Agent the worker should present.
	UserAgent string

	// Proxy is the required proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeVKCaptcha.
func (t *VKCaptchaTask) Type() unicap.TaskType {
	return unicap.TaskTypeVKCaptcha
}

// Validate reports missing required fields and a missing or malformed proxy.
func (t *VKCaptchaTask) Validate() error {
	return validate(
		required("RedirectURI", t.RedirectURI),
		required("UserAgent", t.UserAgent),
		requiredProxy(t.Proxy),
	)
}

// YidunTask solves a NetEase Yidun (NeCaptcha) challenge.
type YidunTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// WebsiteKey is the Yidun captcha ID.
	WebsiteKey string

	// UserAgent is the browser User-Agent the worker should present.
	UserAgent string

	// GetLibURL is the URL of the Yidun loader script, for enterprise
	// deployments.
	GetLibURL string

	// APIServerSubdomain is the custom Yidun API host, for enterprise
	// deployments.
	APIServerSubdomain string

	// Challenge is the enterprise challenge value.
	Challenge string

	// HCG is the enterprise hcg value.
	HCG string

	// HCT is the enterprise hct value.
	HCT int64

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeYidun.
func (t *YidunTask) Type() unicap.TaskType {
	return unicap.TaskTypeYidun
}

// Validate reports missing required fields and malformed proxies.
func (t *YidunTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("WebsiteKey", t.WebsiteKey),
		optionalProxy(t.Proxy),
	)
}

// BinanceTask solves a Binance captcha.
type BinanceTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// WebsiteKey is the Binance bizId value.
	WebsiteKey string

	// ValidateID is the validateId value from the captcha request.
	ValidateID string

	// UserAgent is the browser User-Agent the worker should present.
	UserAgent string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeBinance.
func (t *BinanceTask) Type() unicap.TaskType {
	return unicap.TaskTypeBinance
}

// Validate reports missing required fields and malformed proxies.
func (t *BinanceTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("WebsiteKey", t.WebsiteKey),
		required("ValidateID", t.ValidateID),
		optionalProxy(t.Proxy),
	)
}

// HuntTask solves a Hunt anti-bot challenge.
type HuntTask struct {
	// WebsiteURL is the protected page.
	WebsiteURL string

	// APIGetLib is the URL of the Hunt loader script.
	APIGetLib string

	// UserAgent is the browser User-Agent the worker should present.
	UserAgent string

	// Data is the meta.token value. Leave it empty to receive the X-HD
	// fingerprint first.
	Data string

	// Proxy is the required proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeHunt.
func (t *HuntTask) Type() unicap.TaskType {
	return unicap.TaskTypeHunt
}

// Validate reports missing required fields and a missing or malformed proxy.
func (t *HuntTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("APIGetLib", t.APIGetLib),
		requiredProxy(t.Proxy),
	)
}

// TSPDTask passes an F5 TSPD challenge and returns cookies per domain in the
// raw solution. Use a sticky proxy for the follow-up requests.
type TSPDTask struct {
	// WebsiteURL is the protected page.
	WebsiteURL string

	// TSPDCookie is the TSPD cookie the challenge page set.
	TSPDCookie string

	// HTMLPageBase64 is the base64-encoded challenge page body.
	HTMLPageBase64 string

	// UserAgent is the browser User-Agent the worker should present.
	UserAgent string

	// Proxy is the required proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeTSPD.
func (t *TSPDTask) Type() unicap.TaskType {
	return unicap.TaskTypeTSPD
}

// Validate reports missing required fields and a missing or malformed proxy.
func (t *TSPDTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("TSPDCookie", t.TSPDCookie),
		required("HTMLPageBase64", t.HTMLPageBase64),
		requiredProxy(t.Proxy),
	)
}

// BasiliskTask solves a Basilisk captcha. The response is in the raw
// solution's data.captcha_response field.
type BasiliskTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// WebsiteKey is the Basilisk site key.
	WebsiteKey string

	// UserAgent is the browser User-Agent the worker should present.
	UserAgent string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeBasilisk.
func (t *BasiliskTask) Type() unicap.TaskType {
	return unicap.TaskTypeBasilisk
}

// Validate reports missing required fields and malformed proxies.
func (t *BasiliskTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("WebsiteKey", t.WebsiteKey),
		optionalProxy(t.Proxy),
	)
}

// ImpervaTask passes an Imperva (Incapsula) challenge and returns cookies per
// domain in the raw solution.
type ImpervaTask struct {
	// WebsiteURL is the protected page.
	WebsiteURL string

	// ScriptURL is the URL of the Incapsula script.
	ScriptURL string

	// Cookies are the current Incapsula cookies in "name=value; ..." form.
	Cookies string

	// UserAgent is the browser User-Agent the worker should present.
	UserAgent string

	// Reese84URLEndpoint is the reese84 endpoint, when the site uses one.
	Reese84URLEndpoint string

	// Proxy is the required proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeImperva.
func (t *ImpervaTask) Type() unicap.TaskType {
	return unicap.TaskTypeImperva
}

// Validate reports missing required fields and a missing or malformed proxy.
func (t *ImpervaTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("ScriptURL", t.ScriptURL),
		required("Cookies", t.Cookies),
		requiredProxy(t.Proxy),
	)
}

// AlibabaTask solves an Alibaba Cloud captcha. The tokens are in the raw
// solution's data.tokens field.
type AlibabaTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// SceneID is the captcha SceneId.
	SceneID string

	// Prefix is the captcha prefix value.
	Prefix string

	// UserID is the optional UserId value.
	UserID string

	// UserUserID is the optional UserUserId value.
	UserUserID string

	// VerifyType is the optional verifyType value.
	VerifyType string

	// Region is the optional region value.
	Region string

	// UserCertifyID is the optional UserCertifyId value.
	UserCertifyID string

	// APIGetLib is the optional URL of the captcha loader script.
	APIGetLib string

	// UserAgent is the browser User-Agent the worker should present.
	UserAgent string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeAlibaba.
func (t *AlibabaTask) Type() unicap.TaskType {
	return unicap.TaskTypeAlibaba
}

// Validate reports missing required fields and malformed proxies.
func (t *AlibabaTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("SceneID", t.SceneID),
		required("Prefix", t.Prefix),
		optionalProxy(t.Proxy),
	)
}
