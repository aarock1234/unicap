package anticaptcha

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

type funCaptchaTask struct {
	Type                     string `json:"type"`
	WebsiteURL               string `json:"websiteURL"`
	WebsitePublicKey         string `json:"websitePublicKey"`
	FuncaptchaAPIJSSubdomain string `json:"funcaptchaApiJSSubdomain,omitempty"`
	Data                     string `json:"data,omitempty"`
	UserAgent                string `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

type turnstileTask struct {
	Type        string `json:"type"`
	WebsiteURL  string `json:"websiteURL"`
	WebsiteKey  string `json:"websiteKey"`
	Action      string `json:"action,omitempty"`
	CData       string `json:"cData,omitempty"`
	ChlPageData string `json:"chlPageData,omitempty"`
	solverapi.ProxyFields
}

type geeTestTask struct {
	Type                      string         `json:"type"`
	WebsiteURL                string         `json:"websiteURL"`
	GT                        string         `json:"gt"`
	Challenge                 string         `json:"challenge,omitempty"`
	Version                   int            `json:"version"`
	InitParameters            *geeTestParams `json:"initParameters,omitempty"`
	GeetestAPIServerSubdomain string         `json:"geetestApiServerSubdomain,omitempty"`
	UserAgent                 string         `json:"userAgent,omitempty"`
	solverapi.ProxyFields
}

type geeTestParams struct {
	RiskType string `json:"riskType"`
}

type amazonTask struct {
	Type            string `json:"type"`
	WebsiteURL      string `json:"websiteURL"`
	WebsiteKey      string `json:"websiteKey"`
	WAFType         string `json:"wafType,omitempty"`
	IV              string `json:"iv,omitempty"`
	Context         string `json:"context,omitempty"`
	ChallengeScript string `json:"challengeScript,omitempty"`
	CaptchaScript   string `json:"captchaScript,omitempty"`
	JSAPIScript     string `json:"jsapiScript,omitempty"`
	solverapi.ProxyFields
}

// wafTypeWidget is the wafType of AWS WAF captchas rendered through the
// JavaScript API instead of the page's window.gokuProps.
const wafTypeWidget = "widget"

func mapFunCaptcha(task *tasks.FunCaptchaTask) (funCaptchaTask, error) {
	if err := userAgentWithProxyOnly(task.Proxy, task.UserAgent); err != nil {
		return funCaptchaTask{}, err
	}

	taskType, err := variant("FunCaptchaTask").For(task.Proxy)
	if err != nil {
		return funCaptchaTask{}, err
	}

	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return funCaptchaTask{}, err
	}

	return funCaptchaTask{
		Type:                     taskType,
		WebsiteURL:               task.WebsiteURL,
		WebsitePublicKey:         task.WebsitePublicKey,
		FuncaptchaAPIJSSubdomain: task.APIJSSubdomain,
		Data:                     task.Data,
		UserAgent:                task.UserAgent,
		ProxyFields:              proxy,
	}, nil
}

func mapTurnstile(task *tasks.TurnstileTask) (turnstileTask, error) {
	taskType, err := variant("TurnstileTask").For(task.Proxy)
	if err != nil {
		return turnstileTask{}, err
	}

	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return turnstileTask{}, err
	}

	return turnstileTask{
		Type:        taskType,
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.WebsiteKey,
		Action:      task.Action,
		CData:       task.CData,
		ChlPageData: task.PageData,
		ProxyFields: proxy,
	}, nil
}

func mapGeeTest(task *tasks.GeeTestTask) (geeTestTask, error) {
	if err := userAgentWithProxyOnly(task.Proxy, task.UserAgent); err != nil {
		return geeTestTask{}, err
	}

	taskType, err := variant("GeeTestTask").For(task.Proxy)
	if err != nil {
		return geeTestTask{}, err
	}

	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return geeTestTask{}, err
	}

	return geeTestTask{
		Type:                      taskType,
		WebsiteURL:                task.WebsiteURL,
		GT:                        task.GT,
		Challenge:                 task.Challenge,
		Version:                   3,
		GeetestAPIServerSubdomain: task.APIServerSubdomain,
		UserAgent:                 task.UserAgent,
		ProxyFields:               proxy,
	}, nil
}

// mapGeeTestV4 sends initParameters only when the task names a risk type,
// its single documented key.
func mapGeeTestV4(task *tasks.GeeTestV4Task) (geeTestTask, error) {
	if err := userAgentWithProxyOnly(task.Proxy, task.UserAgent); err != nil {
		return geeTestTask{}, err
	}

	taskType, err := variant("GeeTestTask").For(task.Proxy)
	if err != nil {
		return geeTestTask{}, err
	}

	proxy, err := proxyFields(task.Proxy)
	if err != nil {
		return geeTestTask{}, err
	}

	wire := geeTestTask{
		Type:                      taskType,
		WebsiteURL:                task.WebsiteURL,
		GT:                        task.CaptchaID,
		Version:                   4,
		GeetestAPIServerSubdomain: task.APIServerSubdomain,
		UserAgent:                 task.UserAgent,
		ProxyFields:               proxy,
	}

	if task.RiskType != "" {
		wire.InitParameters = &geeTestParams{
			RiskType: task.RiskType,
		}
	}

	return wire, nil
}

// mapAWSWAF solves a widget integration when the task names its jsapi.js
// script, and the page's gokuProps challenge otherwise.
func mapAWSWAF(task *tasks.AWSWAFTask) (amazonTask, error) {
	mapMode := mapAWSWAFGokuProps
	if task.JSAPIScript != "" {
		mapMode = mapAWSWAFWidget
	}

	wire, err := mapMode(task)
	if err != nil {
		return amazonTask{}, err
	}

	wire.Type, err = variant("AmazonTask").For(task.Proxy)
	if err != nil {
		return amazonTask{}, err
	}

	wire.ProxyFields, err = proxyFields(task.Proxy)
	if err != nil {
		return amazonTask{}, err
	}

	return wire, nil
}

// mapAWSWAFGokuProps maps the fields of gokuProps mode, which solves the
// challenge described by the page's window.gokuProps. The caller fills in
// Type and ProxyFields.
func mapAWSWAFGokuProps(task *tasks.AWSWAFTask) (amazonTask, error) {
	if err := solverapi.Check(
		solverapi.Required(
			solverapi.Setting("Key", task.Key != ""),
			solverapi.Setting("IV", task.IV != ""),
			solverapi.Setting("Context", task.Context != ""),
		),
		solverapi.Unsupported(
			solverapi.Setting("ProblemURL", task.ProblemURL != ""),
			solverapi.Setting("APIKey", task.APIKey != ""),
			solverapi.Setting("ExistingToken", task.ExistingToken != ""),
		),
	); err != nil {
		return amazonTask{}, err
	}

	return amazonTask{
		WebsiteURL:      task.WebsiteURL,
		WebsiteKey:      task.Key,
		IV:              task.IV,
		Context:         task.Context,
		ChallengeScript: task.ChallengeScript,
		CaptchaScript:   task.CaptchaScript,
	}, nil
}

// mapAWSWAFWidget maps the fields of widget mode, which solves a captcha
// rendered by jsapi.js with the key passed to AwsWafCaptcha.renderCaptcha.
// The caller fills in Type and ProxyFields.
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
		WebsiteURL:  task.WebsiteURL,
		WebsiteKey:  task.APIKey,
		WAFType:     wafTypeWidget,
		JSAPIScript: task.JSAPIScript,
	}, nil
}
