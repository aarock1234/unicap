package capmonster

import (
	"encoding/base64"

	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// turnstileTask is CapMonster Cloud's Cloudflare Turnstile task. One wire
// type covers both proxy modes.
type turnstileTask struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	PageAction string `json:"pageAction,omitempty"`
	Data       string `json:"data,omitempty"`
	solverapi.ProxyFields
}

// cloudflareChallengeTask is the cookie variant of CapMonster Cloud's
// Turnstile task, which passes a Cloudflare challenge page and returns the
// cf_clearance cookie.
type cloudflareChallengeTask struct {
	Type               string `json:"type"`
	CloudflareTaskType string `json:"cloudflareTaskType"`
	WebsiteURL         string `json:"websiteURL"`
	WebsiteKey         string `json:"websiteKey"`
	HTMLPageBase64     string `json:"htmlPageBase64"`
	UserAgent          string `json:"userAgent"`
	solverapi.ProxyFields
}

// funCaptchaTask is CapMonster Cloud's FunCaptcha task, which is solved only
// through the caller's proxy.
type funCaptchaTask struct {
	Type                     string `json:"type"`
	WebsiteURL               string `json:"websiteURL"`
	WebsitePublicKey         string `json:"websitePublicKey"`
	Data                     string `json:"data,omitempty"`
	FuncaptchaAPIJSSubdomain string `json:"funcaptchaApiJSSubdomain,omitempty"`
	UserAgent                string `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

// geeTestTask is CapMonster Cloud's GeeTest task, for v3 and v4. The gt
// field holds the v4 captcha ID, and the challenge belongs to v3 only.
type geeTestTask struct {
	Type                      string                 `json:"type"`
	WebsiteURL                string                 `json:"websiteURL"`
	GT                        string                 `json:"gt"`
	Challenge                 string                 `json:"challenge,omitempty"`
	Version                   int                    `json:"version"`
	InitParameters            *geeTestInitParameters `json:"initParameters,omitempty"`
	GeetestAPIServerSubdomain string                 `json:"geetestApiServerSubdomain,omitempty"`
	UserAgent                 string                 `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

// geeTestInitParameters is the initParameters object of a GeeTest v4 task.
type geeTestInitParameters struct {
	RiskType string `json:"riskType"`
}

// amazonTask is CapMonster Cloud's AWS WAF task. The website key is the
// widget's API key or the gokuProps key, and the captcha script is jsapi.js
// or captcha.js, depending on the mode.
type amazonTask struct {
	Type            string `json:"type"`
	WebsiteURL      string `json:"websiteURL"`
	WebsiteKey      string `json:"websiteKey"`
	ChallengeScript string `json:"challengeScript,omitempty"`
	CaptchaScript   string `json:"captchaScript,omitempty"`
	Context         string `json:"context,omitempty"`
	IV              string `json:"iv,omitempty"`
	CookieSolution  bool   `json:"cookieSolution"`
	solverapi.ProxyFields
}

// turnstileTaskType is the wire type of the Turnstile task and of the
// Cloudflare challenge variant selected with cloudflareTaskType.
const turnstileTaskType = "TurnstileTask"

// cloudflareTaskTypeClearance is the cloudflareTaskType of the variant that
// returns the cf_clearance cookie.
const cloudflareTaskTypeClearance = "cf_clearance"

// challengeWebsiteKey fills the website key the cf_clearance variant requires
// but does not use. CapMonster Cloud's documentation says of it: "Any string
// may be passed, for example: xxxxxx".
const challengeWebsiteKey = "xxxxxx"

// geeTestTaskType is the wire type of GeeTest v3 and v4 tasks.
const geeTestTaskType = "GeeTestTask"

// GeeTest versions, as sent in the version field.
const (
	geeTestVersion3 = 3
	geeTestVersion4 = 4
)

// mapTurnstile builds CapMonster Cloud's Cloudflare Turnstile task. PageData
// is rejected because it belongs to the challenge-token variant, which also
// requires a User-Agent the task cannot carry.
func mapTurnstile(task *tasks.TurnstileTask) (turnstileTask, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("PageData", task.PageData != ""),
	); err != nil {
		return turnstileTask{}, err
	}

	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return turnstileTask{}, err
	}

	return turnstileTask{
		Type:        turnstileTaskType,
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		PageAction:  task.Action,
		Data:        task.CData,
		ProxyFields: proxy,
	}, nil
}

// mapCloudflareChallenge builds the cf_clearance variant of CapMonster
// Cloud's Turnstile task. It needs the challenge page, the User-Agent that
// received it, and a proxy with a login and a password.
func mapCloudflareChallenge(task *tasks.CloudflareChallengeTask) (cloudflareChallengeTask, error) {
	hasProxyCredentials := task.Proxy != nil && task.Proxy.Username != "" && task.Proxy.Password != ""

	if err := solverapi.Required(
		solverapi.Setting("HTML", task.HTML != ""),
		solverapi.Setting("UserAgent", task.UserAgent != ""),
		solverapi.Setting("Proxy", task.Proxy != nil),
		solverapi.Setting("Proxy username and password", task.Proxy == nil || hasProxyCredentials),
	); err != nil {
		return cloudflareChallengeTask{}, err
	}

	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return cloudflareChallengeTask{}, err
	}

	return cloudflareChallengeTask{
		Type:               turnstileTaskType,
		CloudflareTaskType: cloudflareTaskTypeClearance,
		WebsiteURL:         task.WebsiteURL,
		WebsiteKey:         challengeWebsiteKey,
		HTMLPageBase64:     base64.StdEncoding.EncodeToString([]byte(task.HTML)),
		UserAgent:          task.UserAgent,
		ProxyFields:        proxy,
	}, nil
}

// mapFunCaptcha builds CapMonster Cloud's FunCaptcha task.
func mapFunCaptcha(task *tasks.FunCaptchaTask) (funCaptchaTask, error) {
	proxy, err := requiredProxyFields(task.Proxy)
	if err != nil {
		return funCaptchaTask{}, err
	}

	return funCaptchaTask{
		Type:                     "FunCaptchaTask",
		WebsiteURL:               task.WebsiteURL,
		WebsitePublicKey:         task.WebsitePublicKey,
		Data:                     task.Data,
		FuncaptchaAPIJSSubdomain: task.APIJSSubdomain,
		UserAgent:                task.UserAgent,
		ProxyFields:              proxy,
	}, nil
}

// mapGeeTest builds CapMonster Cloud's GeeTest v3 task.
func mapGeeTest(task *tasks.GeeTestTask) (geeTestTask, error) {
	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return geeTestTask{}, err
	}

	return geeTestTask{
		Type:                      geeTestTaskType,
		WebsiteURL:                task.WebsiteURL,
		GT:                        task.GT,
		Challenge:                 task.Challenge,
		Version:                   geeTestVersion3,
		GeetestAPIServerSubdomain: task.APIServerSubdomain,
		UserAgent:                 task.UserAgent,
		ProxyFields:               proxy,
	}, nil
}

// mapGeeTestV4 builds CapMonster Cloud's GeeTest v4 task. It sends
// initParameters only when the task names a risk type.
func mapGeeTestV4(task *tasks.GeeTestV4Task) (geeTestTask, error) {
	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return geeTestTask{}, err
	}

	wire := geeTestTask{
		Type:                      geeTestTaskType,
		WebsiteURL:                task.WebsiteURL,
		GT:                        task.CaptchaID,
		Version:                   geeTestVersion4,
		GeetestAPIServerSubdomain: task.APIServerSubdomain,
		UserAgent:                 task.UserAgent,
		ProxyFields:               proxy,
	}

	if task.RiskType != "" {
		wire.InitParameters = &geeTestInitParameters{
			RiskType: task.RiskType,
		}
	}

	return wire, nil
}

// mapAWSWAF builds CapMonster Cloud's AWS WAF task. It solves a widget
// integration when the task names its jsapi.js script, and the page's
// gokuProps challenge otherwise. Either way it asks for the aws-waf-token
// cookie instead of a voucher.
func mapAWSWAF(task *tasks.AWSWAFTask) (amazonTask, error) {
	mapMode := mapAWSWAFGokuProps
	if task.JSAPIScript != "" {
		mapMode = mapAWSWAFWidget
	}

	wire, err := mapMode(task)
	if err != nil {
		return amazonTask{}, err
	}

	wire.ProxyFields, err = proxyFields(task.Proxy)
	if err != nil {
		return amazonTask{}, err
	}

	wire.Type = "AmazonTask"
	wire.CookieSolution = true

	return wire, nil
}

// mapAWSWAFGokuProps maps the fields of gokuProps mode, which solves the
// challenge described by the page's window.gokuProps. The caller fills in
// Type, CookieSolution and ProxyFields.
func mapAWSWAFGokuProps(task *tasks.AWSWAFTask) (amazonTask, error) {
	if err := solverapi.Check(
		solverapi.Required(
			solverapi.Setting("Key", task.Key != ""),
			solverapi.Setting("IV", task.IV != ""),
			solverapi.Setting("Context", task.Context != ""),
			solverapi.Setting("ChallengeScript", task.ChallengeScript != ""),
		),
		solverapi.Unsupported(
			solverapi.Setting("APIKey", task.APIKey != ""),
			solverapi.Setting("ProblemURL", task.ProblemURL != ""),
			solverapi.Setting("ExistingToken", task.ExistingToken != ""),
		),
	); err != nil {
		return amazonTask{}, err
	}

	return amazonTask{
		WebsiteURL:      task.WebsiteURL,
		WebsiteKey:      task.Key,
		ChallengeScript: task.ChallengeScript,
		CaptchaScript:   task.CaptchaScript,
		Context:         task.Context,
		IV:              task.IV,
	}, nil
}

// mapAWSWAFWidget maps the fields of widget mode, which solves a captcha
// rendered by jsapi.js with the key passed to AwsWafCaptcha.renderCaptcha.
// The caller fills in Type, CookieSolution and ProxyFields.
func mapAWSWAFWidget(task *tasks.AWSWAFTask) (amazonTask, error) {
	if err := solverapi.Check(
		solverapi.Required(
			solverapi.Setting("APIKey", task.APIKey != ""),
		),
		solverapi.Unsupported(
			solverapi.Setting("Key", task.Key != ""),
			solverapi.Setting("IV", task.IV != ""),
			solverapi.Setting("Context", task.Context != ""),
			solverapi.Setting("ChallengeScript", task.ChallengeScript != ""),
			solverapi.Setting("CaptchaScript", task.CaptchaScript != ""),
			solverapi.Setting("ProblemURL", task.ProblemURL != ""),
			solverapi.Setting("ExistingToken", task.ExistingToken != ""),
		),
	); err != nil {
		return amazonTask{}, err
	}

	return amazonTask{
		WebsiteURL:    task.WebsiteURL,
		WebsiteKey:    task.APIKey,
		CaptchaScript: task.JSAPIScript,
	}, nil
}
