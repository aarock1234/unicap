package tasks

import (
	"github.com/aarock1234/unicap/v3"
)

// MTCaptchaTask solves an MTCaptcha widget.
type MTCaptchaTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// WebsiteKey is the MTCaptcha site key, starting with "MTPublic-".
	WebsiteKey string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeMTCaptcha.
func (t *MTCaptchaTask) Type() unicap.TaskType {
	return unicap.TaskTypeMTCaptcha
}

// Validate reports missing required fields and malformed proxies.
func (t *MTCaptchaTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("WebsiteKey", t.WebsiteKey),
		optionalProxy(t.Proxy),
	)
}

// FriendlyCaptchaTask solves a Friendly Captcha widget. The token only works
// when the page's own widget script is prevented from running.
type FriendlyCaptchaTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// WebsiteKey is the widget's data-sitekey value.
	WebsiteKey string

	// Version is the widget generation, "v1" or "v2". Empty means v1.
	Version string

	// ModuleScript is the URL of the widget's module script, when known.
	ModuleScript string

	// NomoduleScript is the URL of the widget's nomodule script, when known.
	NomoduleScript string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeFriendlyCaptcha.
func (t *FriendlyCaptchaTask) Type() unicap.TaskType {
	return unicap.TaskTypeFriendlyCaptcha
}

// Validate reports missing required fields and malformed proxies.
func (t *FriendlyCaptchaTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("WebsiteKey", t.WebsiteKey),
		optionalProxy(t.Proxy),
	)
}

// LeminTask solves a Lemin cropped captcha. Decode the result with
// unicap.DecodeSolution[LeminSolution].
type LeminTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// CaptchaID is the Lemin captcha ID, such as "CROPPED_d3d4d56_73ca4008925b4f83a8bed59c2dd0df6d".
	CaptchaID string

	// DivID is the ID of the element that hosts the captcha.
	DivID string

	// APIServerSubdomain is the custom Lemin API host, when the site uses
	// one.
	APIServerSubdomain string

	// UserAgent is the browser User-Agent the worker should present.
	UserAgent string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeLemin.
func (t *LeminTask) Type() unicap.TaskType {
	return unicap.TaskTypeLemin
}

// Validate reports missing required fields and malformed proxies.
func (t *LeminTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("CaptchaID", t.CaptchaID),
		required("DivID", t.DivID),
		optionalProxy(t.Proxy),
	)
}

// LeminSolution is the decoded payload of a LeminTask.
type LeminSolution struct {
	Answer      string `json:"answer"`
	ChallengeID string `json:"challenge_id"`
}

// CutCaptchaTask solves a Cutcaptcha widget.
type CutCaptchaTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// MiseryKey is the CUTCAPTCHA_MISERY_KEY value.
	MiseryKey string

	// APIKey is the data-apikey value.
	APIKey string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeCutCaptcha.
func (t *CutCaptchaTask) Type() unicap.TaskType {
	return unicap.TaskTypeCutCaptcha
}

// Validate reports missing required fields and malformed proxies.
func (t *CutCaptchaTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("MiseryKey", t.MiseryKey),
		required("APIKey", t.APIKey),
		optionalProxy(t.Proxy),
	)
}

// ProsopoTask solves a Prosopo Procaptcha widget.
type ProsopoTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// WebsiteKey is the widget's data-sitekey value.
	WebsiteKey string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeProsopo.
func (t *ProsopoTask) Type() unicap.TaskType {
	return unicap.TaskTypeProsopo
}

// Validate reports missing required fields and malformed proxies.
func (t *ProsopoTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("WebsiteKey", t.WebsiteKey),
		optionalProxy(t.Proxy),
	)
}

// AltchaTask solves an Altcha proof-of-work challenge. Set exactly one of
// ChallengeURL or ChallengeJSON.
type AltchaTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// ChallengeURL is the endpoint that issues the challenge.
	ChallengeURL string

	// ChallengeJSON is the challenge object as a JSON string.
	ChallengeJSON string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeAltcha.
func (t *AltchaTask) Type() unicap.TaskType {
	return unicap.TaskTypeAltcha
}

// Validate reports a missing URL, a missing or doubled challenge source, and
// malformed proxies.
func (t *AltchaTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		exactlyOne("ChallengeURL", t.ChallengeURL, "ChallengeJSON", t.ChallengeJSON),
		optionalProxy(t.Proxy),
	)
}

// CapyTask solves a Capy puzzle captcha. Decode the result with
// unicap.DecodeSolution[CapySolution].
type CapyTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// WebsiteKey is the Capy captcha key.
	WebsiteKey string

	// UserAgent is the browser User-Agent the worker should present.
	UserAgent string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeCapy.
func (t *CapyTask) Type() unicap.TaskType {
	return unicap.TaskTypeCapy
}

// Validate reports missing required fields and malformed proxies.
func (t *CapyTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("WebsiteKey", t.WebsiteKey),
		optionalProxy(t.Proxy),
	)
}

// CapySolution is the decoded payload of a CapyTask.
type CapySolution struct {
	CaptchaKey   string `json:"captchakey"`
	ChallengeKey string `json:"challengekey"`
	Answer       string `json:"answer"`
	RespKey      string `json:"respKey"`
}

// KeyCaptchaTask solves a KeyCAPTCHA puzzle.
type KeyCaptchaTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// UserID is the s_s_c_user_id value.
	UserID string

	// SessionID is the s_s_c_session_id value.
	SessionID string

	// WebServerSign is the s_s_c_web_server_sign value.
	WebServerSign string

	// WebServerSign2 is the s_s_c_web_server_sign2 value.
	WebServerSign2 string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeKeyCaptcha.
func (t *KeyCaptchaTask) Type() unicap.TaskType {
	return unicap.TaskTypeKeyCaptcha
}

// Validate reports missing required fields and malformed proxies.
func (t *KeyCaptchaTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("UserID", t.UserID),
		required("SessionID", t.SessionID),
		required("WebServerSign", t.WebServerSign),
		required("WebServerSign2", t.WebServerSign2),
		optionalProxy(t.Proxy),
	)
}

// CyberSiAraTask solves a CyberSiARA slider.
type CyberSiAraTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// SlideMasterURLID is the MasterUrlId value from the widget.
	SlideMasterURLID string

	// UserAgent is the browser User-Agent the worker should present.
	UserAgent string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeCyberSiAra.
func (t *CyberSiAraTask) Type() unicap.TaskType {
	return unicap.TaskTypeCyberSiAra
}

// Validate reports missing required fields and malformed proxies.
func (t *CyberSiAraTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("SlideMasterURLID", t.SlideMasterURLID),
		required("UserAgent", t.UserAgent),
		optionalProxy(t.Proxy),
	)
}

// AtbCaptchaTask solves an atbCAPTCHA widget.
type AtbCaptchaTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// AppID is the widget's appId value.
	AppID string

	// APIServer is the widget's apiServer value.
	APIServer string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeAtbCaptcha.
func (t *AtbCaptchaTask) Type() unicap.TaskType {
	return unicap.TaskTypeAtbCaptcha
}

// Validate reports missing required fields and malformed proxies.
func (t *AtbCaptchaTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("AppID", t.AppID),
		required("APIServer", t.APIServer),
		optionalProxy(t.Proxy),
	)
}

// TencentTask solves a Tencent (TenDI) captcha. Decode the result with
// unicap.DecodeSolution[TencentSolution].
type TencentTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// AppID is the captcha's appId value.
	AppID string

	// CaptchaScript is the URL of the captcha script, when the site loads a
	// custom one.
	CaptchaScript string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeTencent.
func (t *TencentTask) Type() unicap.TaskType {
	return unicap.TaskTypeTencent
}

// Validate reports missing required fields and malformed proxies.
func (t *TencentTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("AppID", t.AppID),
		optionalProxy(t.Proxy),
	)
}

// TencentSolution is the decoded payload of a TencentTask.
type TencentSolution struct {
	AppID   string `json:"appid"`
	Ret     int    `json:"ret"`
	Ticket  string `json:"ticket"`
	RandStr string `json:"randstr"`
}

// CaptchaFoxTask solves a CaptchaFox widget.
type CaptchaFoxTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// WebsiteKey is the widget's site key.
	WebsiteKey string

	// UserAgent is the browser User-Agent the worker should present.
	UserAgent string

	// APIServer is the custom CaptchaFox API host, when the site uses one.
	APIServer string

	// Proxy is the required proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeCaptchaFox.
func (t *CaptchaFoxTask) Type() unicap.TaskType {
	return unicap.TaskTypeCaptchaFox
}

// Validate reports missing required fields and a missing or malformed proxy.
func (t *CaptchaFoxTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("WebsiteKey", t.WebsiteKey),
		required("UserAgent", t.UserAgent),
		requiredProxy(t.Proxy),
	)
}

// YandexSmartCaptchaTask solves a Yandex SmartCaptcha widget. Submit the
// token in the page's smart-token field.
type YandexSmartCaptchaTask struct {
	// WebsiteURL is the page that shows the captcha.
	WebsiteURL string

	// WebsiteKey is the widget's data-sitekey value.
	WebsiteKey string

	// UserAgent is the browser User-Agent the worker should present.
	UserAgent string

	// Cookies are page cookies in "name=value; name2=value2" form.
	Cookies string

	// Proxy is the optional proxy the worker loads the page through.
	Proxy *unicap.Proxy
}

// Type returns unicap.TaskTypeYandexSmartCaptcha.
func (t *YandexSmartCaptchaTask) Type() unicap.TaskType {
	return unicap.TaskTypeYandexSmartCaptcha
}

// Validate reports missing required fields and malformed proxies.
func (t *YandexSmartCaptchaTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("WebsiteKey", t.WebsiteKey),
		optionalProxy(t.Proxy),
	)
}

// BotDeflectorTask solves a BotDeflector challenge and returns a JWT in
// unicap.Solution.Token.
type BotDeflectorTask struct {
	// WebsiteURL is the protected page.
	WebsiteURL string

	// Domain is the BotDeflector domain value.
	Domain string

	// FlowToken is the BotDeflector flow token.
	FlowToken string
}

// Type returns unicap.TaskTypeBotDeflector.
func (t *BotDeflectorTask) Type() unicap.TaskType {
	return unicap.TaskTypeBotDeflector
}

// Validate reports missing required fields.
func (t *BotDeflectorTask) Validate() error {
	return validate(
		required("WebsiteURL", t.WebsiteURL),
		required("Domain", t.Domain),
		required("FlowToken", t.FlowToken),
	)
}
